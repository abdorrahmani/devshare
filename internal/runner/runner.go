package runner

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/abdorrahmani/devshare/internal/middleware"
	"github.com/abdorrahmani/devshare/internal/network"
	"github.com/abdorrahmani/devshare/internal/qrcode"
)

// RunProject runs the appropriate command based on project type and package manager.
// It supports Laravel, React, Next.js, Go, and Node.js projects.
func RunProject(projectType, packageManager, port, password string) error {
	switch projectType {
	case "laravel":
		return runLaravel(port, password)
	case "react":
		return runReact(packageManager, port, password)
	case "nextjs":
		return runNextJS(packageManager, port, password)
	case "go":
		return runGo(password)
	case "nodejs":
		return runNodeJS(packageManager, port, password)
	default:
		return fmt.Errorf("unsupported project type: %s", projectType)
	}
}

// startAuthServer starts an authentication proxy server that forwards requests to the actual app
func startAuthServer(targetPort, authPort, password string) error {
	ip := network.GetLocalIP()
	if ip == "" {
		return fmt.Errorf("could not determine local IP address")
	}

	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		targetURL := fmt.Sprintf("http://localhost:%s%s", targetPort, r.URL.RequestURI())
		proxyReq, err := http.NewRequest(r.Method, targetURL, io.NopCloser(bytes.NewBuffer(bodyBytes)))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for name, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}

		proxyReq.Header.Set("Host", fmt.Sprintf("localhost:%s", targetPort))

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		io.Copy(w, resp.Body)
	})

	var handler http.Handler = proxyHandler
	if password != "" {
		handler = middleware.AuthMiddleware(proxyHandler, password)
		fmt.Printf("🔐 Authentication enabled - Password required to access the app\n")
	}

	server := &http.Server{
		Addr:    ":" + authPort,
		Handler: handler,
	}

	fmt.Printf("🔗 Auth Proxy: http://%s:%s\n", ip, authPort)
	qrcode.GenerateQrCodeWithMessage(ip+":"+authPort, "📱 Scan this on your phone:")

	return server.ListenAndServe()
}

// runWithInstallRetry tries to run the app with given commands, installs dependencies if needed, and retries. Shows QR code only after successful start.
func runWithInstallRetry(packageManager string, cmds [][]string, installArgs []string, port, password string) error {
	if password != "" {
		authPort := strconv.Itoa(getAvailablePort(port))
		go func() {
			if err := startAuthServer(port, authPort, password); err != nil {
				fmt.Printf("❌ Auth server error: %v\n", err)
			}
		}()

		time.Sleep(500 * time.Millisecond)
	}

	for _, args := range cmds {
		cmd := exec.Command(packageManager, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err == nil {
			if password == "" {
				ip := network.GetLocalIP()
				qrcode.GenerateQrCodeWithMessage(ip+":"+port, "📱 Scan this on your phone:")
			}
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
			if password == "" {
				ip := network.GetLocalIP()
				qrcode.GenerateQrCodeWithMessage(ip+":"+port, "📱 Scan this on your phone:")
			}
			cmd.Wait()
			return nil
		}
	}
	return fmt.Errorf("failed to start app with %s", packageManager)
}

// getAvailablePort finds an available port starting from the given port
func getAvailablePort(startPort string) int {
	port, _ := strconv.Atoi(startPort)
	if port == 0 {
		port = 3000
	}
	startSearchPort := port + 1

	for i := 0; i < 100; i++ {
		testPort := startSearchPort + i
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(testPort))
		if err == nil {
			listener.Close()
			return testPort
		}
	}
	return startSearchPort + 1
}

func runReact(packageManager, port, password string) error {
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
	if password == "" {
		fmt.Printf("Network: http://%s:%s\n", ip, port)
	}
	return runWithInstallRetry(
		packageManager,
		cmds,
		[]string{"install"},
		port,
		password,
	)
}

func runNextJS(packageManager, port, password string) error {
	fmt.Println("🚀 Starting Next.js app...")
	ip := network.GetLocalIP()
	if port == "" {
		port = "3000" // Default Next.js port
	}
	cmds := [][]string{
		{"dev", "--port", port, "-H", "0.0.0.0"},
	}
	fmt.Printf("Local:   http://localhost:%s\n", port)
	if password == "" {
		fmt.Printf("Network: http://%s:%s\n", ip, port)
	}
	return runWithInstallRetry(
		packageManager,
		cmds,
		[]string{"install"},
		port,
		password,
	)
}

func runGo(password string) error {
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
func runLaravel(port, password string) error {
	fmt.Println("🚀 Starting Laravel app...")
	ip := network.GetLocalIP()
	if port == "" {
		port = "8000" // Default Laravel port
	}
	cmd := exec.Command("php", "artisan", "serve", "--host", "0.0.0.0", "--port", port)

	fmt.Printf("Local:   http://localhost:%s\n", port)
	if password == "" {
		fmt.Printf("Network: http://%s:%s\n", ip, port)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if password != "" {
		authPort := strconv.Itoa(getAvailablePort(port))
		go func() {
			if err := startAuthServer(port, authPort, password); err != nil {
				fmt.Printf("❌ Auth server error: %v\n", err)
			}
		}()
		time.Sleep(500 * time.Millisecond)
	} else {
		qrcode.GenerateQrCodeWithMessage(ip+":"+port, "📱 Scan this on your phone:")
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Laravel app: %w", err)
	}
	return nil
}

// runNodeJs runs the Node.js app
func runNodeJS(packageManager, port, password string) error {
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
		if password == "" {
			fmt.Printf("Network: http://%s:%s\n", ip, port)
			qrcode.GenerateQrCodeWithMessage(ip+":"+port, "📱 Scan this on your phone:")
		}
	}

	if password != "" {
		authPort := strconv.Itoa(getAvailablePort(port))
		go func() {
			if err := startAuthServer(port, authPort, password); err != nil {
				fmt.Printf("❌ Auth server error: %v\n", err)
			}
		}()
		time.Sleep(500 * time.Millisecond)
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
