package cmd

import (
	"github.com/abdorrahmani/devshare/internal/display"
	"github.com/spf13/cobra"
)

var (
	Version = "v0.1.0"
)

var rootCmd = &cobra.Command{
	Use:   "devshare",
	Short: "DevShare is a CLI tool for sharing your dev environment over LAN",
	Long: `DevShare is a CLI tool for sharing your development environment over LAN.

Example usage:
  devshare
  devshare --version
`,
	Run: func(cmd *cobra.Command, args []string) {
		display.WelcomeMessage()
	},
}

func init() {
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("DevShare version: {{.Version}}\n")
}

func Execute() error {
	return rootCmd.Execute()
}
