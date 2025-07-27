package runner

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/abdorrahmani/devshare/internal/network"

	"github.com/abdorrahmani/devshare/internal/qrcode"
)

// RunProject runs the appropriate command based on project type and package manager.
// It supports Laravel, React, Next.js, Go, and Node.js projects.
func RunProject(projectType, packageManager, port string) error {
	switch projectType {
	case "laravel":
		return runLaravel(port)
	case "react":
		return runReact(packageManager, port)
	case "nextjs":
		return runNextJS(packageManager, port)
	case "go":
		return runGo()
	case "nodejs":
		return runNodeJS(packageManager, port)
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

func runReact(packageManager, port string) error {
	fmt.Println("🚀 Starting React app...")
	ip := network.GetLocalIP()
	if port == "" {
		port = "5173" // Default Vite port
	}
	cmds := [][]string{
		{"start", "--port", port, "--host", "0.0.0.0"},
		{"dev", "--port", port, "--host", "0.0.0.0"},
	}
	fmt.Printf("Local:   http://localhost:%s\n", port)
	fmt.Printf("Network: http://%s:%s\n", ip, port)
	return runWithInstallRetry(
		packageManager,
		cmds,
		[]string{"install"},
		ip+":"+port,
		"📱 Scan this on your phone:",
	)
}

func runNextJS(packageManager, port string) error {
	fmt.Println("🚀 Starting Next.js app...")
	ip := network.GetLocalIP()
	if port == "" {
		port = "3000" // Default Next.js port
	}
	cmds := [][]string{
		{"dev", "--port", port, "-H", "0.0.0.0"},
	}
	fmt.Printf("Local:   http://localhost:%s\n", port)
	fmt.Printf("Network: http://%s:%s\n", ip, port)
	return runWithInstallRetry(
		packageManager,
		cmds,
		[]string{"install"},
		ip+":"+port,
		"📱 Scan this on your phone:",
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

// runLaravel runs the Laravel app
func runLaravel(port string) error {
	fmt.Println("🚀 Starting Laravel app...")
	ip := network.GetLocalIP()
	if port == "" {
		port = "8000" // Default Laravel port
	}
	cmd := exec.Command("php", "artisan", "serve", "--host", "0.0.0.0", "--port", port)

	fmt.Printf("Local:   http://localhost:%s\n", port)
	fmt.Printf("Network: http://%s:%s\n", ip, port)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	qrcode.GenerateQrCodeWithMessage(ip+":"+port, "📱 Scan this on your phone:")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Laravel app: %w", err)
	}
	return nil
}

// runNodeJs runs the Node.js app
func runNodeJS(packageManager, port string) error {
	fmt.Println("🚀 Starting Node.js app...")
	ip := network.GetLocalIP()

	if port == "" {
		port = "3000" // Default Node.js port
	}
	pmCmds := [][]string{
		{"start"},
		{"run", "dev"},
	}

	entryFiles := []struct {
		file  string
		useTs bool
	}{
		{"index.js", false},
		{"app.js", false},
		{"index.ts", true},
		{"app.ts", true},
	}

	printNetworkInfo := func() {
		fmt.Printf("Local:   http://localhost:%s\n", port)
		fmt.Printf("Network: http://%s:%s\n", ip, port)
		qrcode.GenerateQrCodeWithMessage(ip+":"+port, "📱 Scan this on your phone:")
	}

	// Try package manager scripts first
	for _, args := range pmCmds {
		fmt.Printf("Trying: %s %s\n", packageManager, args)
		printNetworkInfo()
		cmd := exec.Command(packageManager, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		} else {
			fmt.Printf("⚠️  Failed to run %s %s: %v\n", packageManager, args, err)
		}
	}

	// Try direct node/ts-node with entry files
	for _, entry := range entryFiles {
		if _, err := os.Stat(entry.file); err == nil {
			var cmd *exec.Cmd
			if entry.useTs {
				fmt.Printf("Trying: ts-node %s\n", entry.file)
				printNetworkInfo()
				cmd = exec.Command("ts-node", entry.file)
			} else {
				fmt.Printf("Trying: node %s\n", entry.file)
				printNetworkInfo()
				cmd = exec.Command("node", entry.file)
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err == nil {
				return nil
			} else {
				fmt.Printf("⚠️  Failed to run %s %s: %v\n", cmd.Path, entry.file, err)
			}
		}
	}

	fmt.Println("Installing dependencies...")
	installCmd := exec.Command(packageManager, "install")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	// Retry package manager scripts
	for _, args := range pmCmds {
		fmt.Printf("Trying: %s %s\n", packageManager, args)
		printNetworkInfo()
		cmd := exec.Command(packageManager, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		} else {
			fmt.Printf("⚠️  Failed to run %s %s: %v\n", packageManager, args, err)
		}
	}

	// Retry direct node/ts-node
	for _, entry := range entryFiles {
		if _, err := os.Stat(entry.file); err == nil {
			var cmd *exec.Cmd
			if entry.useTs {
				fmt.Printf("Trying: ts-node %s\n", entry.file)
				printNetworkInfo()
				cmd = exec.Command("ts-node", entry.file)
			} else {
				fmt.Printf("Trying: node %s\n", entry.file)
				printNetworkInfo()
				cmd = exec.Command("node", entry.file)
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err == nil {
				return nil
			} else {
				fmt.Printf("⚠️  Failed to run %s %s: %v\n", cmd.Path, entry.file, err)
			}
		}
	}

	return fmt.Errorf("could not start Node.js app: no working package manager script or entry file")
}
