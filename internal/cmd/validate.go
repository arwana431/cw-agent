package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/certwatch-app/cw-agent/internal/config"
	"github.com/certwatch-app/cw-agent/internal/ui"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file",
	Long: `Validate the CertWatch Agent configuration file without starting the agent.

Example:
  cw-agent validate -c /path/to/certwatch.yaml`,
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Println(ui.RenderCommandHeader("Config Validation"))
	fmt.Println()

	// Load configuration
	cfg, err := config.Load(viper.GetViper())
	if err != nil {
		fmt.Println(ui.RenderError("Failed to load configuration"))
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	fmt.Println(ui.RenderSuccess("Configuration loaded"))

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Println(ui.RenderError("Configuration validation failed"))
		fmt.Println(ui.RenderError("  " + err.Error()))
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	fmt.Println(ui.RenderSuccess("API settings valid"))
	fmt.Println(ui.RenderSuccess("Agent settings valid"))
	fmt.Println(ui.RenderSuccess(fmt.Sprintf("%d certificates configured", len(cfg.Certificates))))

	// Summary section
	fmt.Println()
	fmt.Println(ui.RenderSection("Summary"))
	fmt.Println()
	fmt.Println(ui.RenderKeyValue("Agent", cfg.Agent.Name))

	// Only show API endpoint if non-default
	if cfg.API.Endpoint != ui.DefaultAPIEndpoint {
		fmt.Println(ui.RenderKeyValue("API Endpoint", cfg.API.Endpoint))
	}

	fmt.Println(ui.RenderKeyValue("Certificates", fmt.Sprintf("%d", len(cfg.Certificates))))
	fmt.Println(ui.RenderKeyValue("Sync", cfg.Agent.SyncInterval.String()))
	fmt.Println(ui.RenderKeyValue("Scan", cfg.Agent.ScanInterval.String()))
	fmt.Println()

	fmt.Println(ui.RenderSuccess("Configuration is valid!"))
	fmt.Println()

	return nil
}
