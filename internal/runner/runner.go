package runner

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/abdorrahmani/devshare/internal/network"

	"github.com/abdorrahmani/devshare/internal/qrcode"
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

// runWithInstallRetry tries to run the app with given commands, installs dependencies if needed, and retries. Shows QR code only after successful start.
func runWithInstallRetry(packageManager string, cmds [][]string, installArgs []string, qrUrl string, qrMsg string) error {
	for _, args := range cmds {
		cmd := exec.Command(packageManager, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err == nil {
			qrcode.GenerateQrCodeWithMessage(qrUrl, qrMsg)
			cmd.Wait()
			return nil
		}
	}
	fmt.Println("Installing dependencies...")
	installCmd := exec.Command(packageManager, installArgs...)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}
	for _, args := range cmds {
		cmd := exec.Command(packageManager, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err == nil {
			qrcode.GenerateQrCodeWithMessage(qrUrl, qrMsg)
			cmd.Wait()
			return nil
		}
	}
	return fmt.Errorf("failed to start app with %s", packageManager)
}

func runReact(packageManager string) error {
	fmt.Println("🚀 Starting React app...")
	ip := network.GetLocalIP()
	cmds := [][]string{
		{"start", "--host", "0.0.0.0"},
		{"dev", "--host", "0.0.0.0"},
	}
	fmt.Printf("Local:   http://localhost:5173\n")
	fmt.Printf("Network: http://%s:5173\n", ip)
	return runWithInstallRetry(
		packageManager,
		cmds,
		[]string{"install"},
		ip+":5173",
		"📱 Scan this on your phone (React/Vite default port 5173):",
	)
}

func runNextJS(packageManager string) error {
	fmt.Println("🚀 Starting Next.js app...")
	ip := network.GetLocalIP()
	cmds := [][]string{
		{"dev", "-H", "0.0.0.0"},
	}
	fmt.Printf("Local:   http://localhost:3000\n")
	fmt.Printf("Network: http://%s:3000\n", ip)
	return runWithInstallRetry(
		packageManager,
		cmds,
		[]string{"install"},
		ip+":3000",
		"📱 Scan this on your phone (Next.js default port 3000):",
	)
}

func runGo() error {
	fmt.Println("🚀 Starting Go app...")
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Go app: %w", err)
	}
	return nil
}
