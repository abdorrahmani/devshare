package runner

import (
	"fmt"
	"os"
	"os/exec"
)

// RunProject runs the appropriate command based on project type and package manager.
func RunProject(projectType, packageManager string) error {
	switch projectType {
	case "react":
		return runReact(packageManager)
	case "nextjs":
		return runNextJS(packageManager)
	case "go":
		return runGo()
	default:
		return fmt.Errorf("unsupported project type: %s", projectType)
	}
}

func runReact(packageManager string) error {
	fmt.Println("🚀 Starting React app...")
	cmds := [][]string{
		{"start", "--host", "0.0.0.0"},
		{"dev", "--host", "0.0.0.0"},
	}

	for _, args := range cmds {
		cmd := exec.Command(packageManager, args...)
		cmd.Stderr = os.Stdout
		cmd.Stdout = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try to install then retry
	fmt.Println("Installing dependencies...")
	installCmd := exec.Command(packageManager, "install")
	installCmd.Stdout = os.Stdout
	installCmd.Stdout = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}
	// Retry start
	for _, args := range cmds {
		cmd := exec.Command(packageManager, args...)
		cmd.Stdout = os.Stdout
		cmd.Stdout = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to start React app with %s", packageManager)
}

func runNextJS(packageManager string) error {
	fmt.Println("🚀 Starting Next.js app...")
	cmds := [][]string{
		{"dev", "-H", "0.0.0.0"},
	}
	for _, args := range cmds {
		cmd := exec.Command(packageManager, args...)
		cmd.Stdout = os.Stdout
		cmd.Stdout = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	fmt.Println("Installing dependencies...")
	installCmd := exec.Command(packageManager, "install")
	installCmd.Stdout = os.Stdout
	installCmd.Stdout = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}
	for _, args := range cmds {
		cmd := exec.Command(packageManager, args...)
		cmd.Stdout = os.Stdout
		cmd.Stdout = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to start Next.js app with %s", packageManager)
}

func runGo() error {
	fmt.Println("🚀 Starting Go app...")
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stdout = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Go app: %w", err)
	}
	return nil
}
