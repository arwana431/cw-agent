package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/certwatch-app/cw-agent/internal/agent"
	"github.com/certwatch-app/cw-agent/internal/config"
	"github.com/certwatch-app/cw-agent/internal/state"
	"github.com/certwatch-app/cw-agent/internal/ui"
)

var (
	resetAgent bool
	yesFlag    bool
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the CertWatch monitoring agent",
	Long: `Start the CertWatch Agent to monitor configured certificates
and sync data to the CertWatch cloud platform.

Example:
  cw-agent start -c /path/to/certwatch.yaml
  cw-agent start --config certwatch.yaml
  cw-agent start -c certwatch.yaml --reset-agent --yes`,
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVar(&resetAgent, "reset-agent", false,
		"Reset agent state and re-register (transfers certs to new agent, orphans removed certs)")
	startCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false,
		"Skip confirmation prompts (for CI/automation)")
}

func runStart(cmd *cobra.Command, args []string) error {
	// Load and validate configuration
	cfg, err := config.Load(viper.GetViper())
	if err != nil {
		fmt.Println()
		fmt.Println(ui.RenderError("Failed to load configuration: " + err.Error()))
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if validationErr := cfg.Validate(); validationErr != nil {
		fmt.Println()
		fmt.Println(ui.RenderError("Invalid configuration: " + validationErr.Error()))
		return fmt.Errorf("invalid configuration: %w", validationErr)
	}

	// Initialize state manager using the config file path from viper
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		configPath = "./certwatch.yaml" // fallback
	}
	stateManager := state.NewManager(configPath)
	if loadErr := stateManager.Load(); loadErr != nil {
		// Log warning but continue - corrupted state is treated as first run
		fmt.Println(ui.RenderWarning(loadErr.Error()))
	}

	// Check for name change (only if we have existing state and not using --reset-agent)
	if !resetAgent && stateManager.HasNameChanged(cfg.Agent.Name) {
		return handleNameChangeWarning(stateManager, cfg)
	}

	// Handle reset flag
	if resetAgent && stateManager.HasState() {
		if resetErr := handleAgentReset(stateManager, cfg); resetErr != nil {
			return resetErr
		}
	}

	// Create agent with state manager
	a, err := agent.New(cfg, stateManager)
	if err != nil {
		fmt.Println(ui.RenderError("Failed to create agent: " + err.Error()))
		return fmt.Errorf("failed to create agent: %w", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Println()
		fmt.Println(ui.RenderWarning(fmt.Sprintf("Received signal %v, shutting down...", sig)))
		cancel()
	}()

	// Display styled startup info
	fmt.Println()
	fmt.Println(ui.RenderAppHeader())

	fmt.Println(ui.RenderSection("Agent"))
	fmt.Println()
	fmt.Println(ui.RenderKeyValue("Name", cfg.Agent.Name))

	if agentID := stateManager.GetAgentID(); agentID != "" {
		fmt.Println(ui.RenderKeyValue("Agent ID", ui.TruncateID(agentID)))
	}

	fmt.Println(ui.RenderKeyValue("Certificates", fmt.Sprintf("%d", len(cfg.Certificates))))
	fmt.Println(ui.RenderKeyValue("Sync", cfg.Agent.SyncInterval.String()))
	fmt.Println()

	fmt.Println(ui.RenderSuccess("Agent started"))
	fmt.Println()

	// Run the agent
	if err := a.Run(ctx); err != nil && err != context.Canceled {
		fmt.Println(ui.RenderError("Agent error: " + err.Error()))
		return fmt.Errorf("agent error: %w", err)
	}

	fmt.Println()
	fmt.Println(ui.RenderInfo("Agent stopped gracefully"))
	return nil
}

// handleNameChangeWarning displays a warning when agent name has changed and exits
func handleNameChangeWarning(sm *state.Manager, cfg *config.Config) error {
	fmt.Println()
	fmt.Println(ui.RenderAppHeader())
	fmt.Println()

	// Build warning box content
	lines := []string{
		"",
		fmt.Sprintf("  Previous: %q (%s)", sm.GetAgentName(), ui.TruncateID(sm.GetAgentID())),
		fmt.Sprintf("  New:      %q", cfg.Agent.Name),
		"",
		"  Certificates from the old agent will NOT be",
		"  automatically transferred.",
		"",
		"  Options:",
		"  1. Continue with new name (migrates matching certs):",
		"     " + ui.RenderCode("cw-agent start -c <config> --reset-agent"),
		"",
		"  2. Keep existing agent (revert config):",
		fmt.Sprintf("     Edit config: agent.name = %q", sm.GetAgentName()),
	}
	fmt.Println(ui.RenderWarningBox("Agent Name Changed", lines))
	fmt.Println()

	fmt.Println(ui.RenderInfo("Exiting. No changes made."))
	fmt.Println()

	return fmt.Errorf("agent name changed - use --reset-agent to continue with new name")
}

// handleAgentReset handles the --reset-agent flag, asking for confirmation unless --yes is provided
func handleAgentReset(sm *state.Manager, cfg *config.Config) error {
	fmt.Println()
	fmt.Println(ui.RenderAppHeader())
	fmt.Println()

	// Build warning box content
	lines := []string{
		"",
		fmt.Sprintf("  Previous: %q (%s)", sm.GetAgentName(), ui.TruncateID(sm.GetAgentID())),
		fmt.Sprintf("  New:      %q", cfg.Agent.Name),
		"",
		"  On next sync:",
		"  • Matching certs will migrate to new agent",
		"  • Non-matching certs will be orphaned",
	}
	fmt.Println(ui.RenderWarningBox("Agent Reset", lines))
	fmt.Println()

	// Skip confirmation if --yes flag is provided
	if !yesFlag {
		if !promptConfirm("Continue?") {
			fmt.Println()
			fmt.Println(ui.RenderWarning("Reset canceled by user"))
			return fmt.Errorf("reset canceled by user")
		}
	} else {
		fmt.Println(ui.RenderInfo("--yes flag provided, skipping confirmation"))
	}

	// Store previous agent ID for migration
	previousAgentID := sm.GetAgentID()
	if previousAgentID != "" {
		sm.SetPreviousAgentID(previousAgentID)
	}

	// Clear current agent ID (will be assigned on next sync)
	sm.SetAgentID("")

	// Save the state with previous_agent_id
	if err := sm.Save(); err != nil {
		fmt.Println(ui.RenderError("Failed to save state: " + err.Error()))
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Println()
	fmt.Println(ui.RenderSuccess("Agent state reset"))
	fmt.Println(ui.RenderInfo("Starting with new agent..."))
	fmt.Println()

	return nil
}

// promptConfirm asks for user confirmation
func promptConfirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", prompt)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
