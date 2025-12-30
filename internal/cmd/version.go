package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/certwatch-app/cw-agent/internal/ui"
	"github.com/certwatch-app/cw-agent/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version, git commit, and build date of the CertWatch Agent.`,
	Run: func(cmd *cobra.Command, args []string) {
		info := version.GetInfo()

		fmt.Println()
		fmt.Println(ui.RenderAppHeader())
		fmt.Println()

		fmt.Println(ui.RenderKeyValue("Version", info.Version))
		fmt.Println(ui.RenderKeyValue("Commit", info.GitCommit))
		fmt.Println(ui.RenderKeyValue("Build Date", info.BuildDate))
		fmt.Println(ui.RenderKeyValue("Go Version", info.GoVersion))
		fmt.Println(ui.RenderKeyValue("Platform", info.OS+"/"+info.Arch))
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
