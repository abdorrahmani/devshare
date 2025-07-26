package cmd

import (
	"fmt"

	"github.com/abdorrahmani/devshare/internal/detector"
	"github.com/abdorrahmani/devshare/internal/display"
	"github.com/abdorrahmani/devshare/internal/runner"
	"github.com/spf13/cobra"
)

var (
	Version = "v1.0.0"
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
		dir := detector.GetWorkingDir()
		projectType, pkgManager := detector.DetectProjectType(dir)
		if projectType == "" {
			fmt.Println("❌ No supported project detected. Exiting.")
			return
		}
		if projectType == "react" || projectType == "nextjs" || projectType == "nodejs" {
			err := runner.RunProject(projectType, pkgManager)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}
		} else if projectType == "go" || projectType == "laravel" {
			err := runner.RunProject(projectType, "")
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}
		} else {
			fmt.Printf("Project type '%s' detected, but no runner implemented.\n", projectType)
		}
	},
}

func init() {
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("DevShare version: {{.Version}}\n")
}

func Execute() error {
	return rootCmd.Execute()
}
