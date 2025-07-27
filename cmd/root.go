package cmd

import (
	"flag"
	"fmt"

	"github.com/abdorrahmani/devshare/internal/detector"
	"github.com/abdorrahmani/devshare/internal/display"
	"github.com/abdorrahmani/devshare/internal/runner"
	"github.com/spf13/cobra"
)

var (
	Version      = "v1.1.0"
	passwordFlag = flag.String("password", "", "Password for accessing the shared environment (optional, defaults to no password)")
)

var rootCmd = &cobra.Command{
	Use:   "devshare [PORT] --password [PASSWORD]",
	Args:  cobra.MaximumNArgs(1),
	Short: "DevShare is a CLI tool for sharing your dev environment over LAN",
	Long: `DevShare is a CLI tool for sharing your development environment over LAN.

Example usage:
  devshare
  devshare 3000
  devshare --version
`,
	Run: func(cmd *cobra.Command, args []string) {
		display.WelcomeMessage()

		var port string
		if len(args) > 0 {
			port = args[0]
		} else {
			port = ""
		}
		fmt.Println("Password:", *passwordFlag)
		dir := detector.GetWorkingDir()
		projectType, pkgManager := detector.DetectProjectType(dir)
		if projectType == "" {
			fmt.Println("❌ No supported project detected. Exiting.")
			return
		}
		if projectType == "react" || projectType == "nextjs" || projectType == "nodejs" {
			err := runner.RunProject(projectType, pkgManager, port, *passwordFlag)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}
		} else if projectType == "go" || projectType == "laravel" {
			err := runner.RunProject(projectType, "", port, *passwordFlag)
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
	rootCmd.Flags().StringVar(passwordFlag, "password", "", "Password for accessing the shared environment (optional, defaults to no password)")
	rootCmd.SetVersionTemplate("DevShare version: {{.Version}}\n")
}

func Execute() error {
	return rootCmd.Execute()
}
