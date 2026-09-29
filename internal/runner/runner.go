package runner

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/abdorrahmani/devshare/internal/middleware"
	"github.com/abdorrahmani/devshare/internal/network"
	"github.com/abdorrahmani/devshare/internal/qrcode"
)

// RunProject runs the appropriate command based on project type and package manager.
// It supports Laravel, React, Vue, Next.js, Go, and Node.js projects.
func RunProject(projectType, packageManager, port, password string) error {
	switch projectType {
	case "laravel":
		return runLaravel(port, password)
	case "react":
		return runVite("React", packageManager, port, password)
	case "vue":
		return runVite("Vue", packageManager, port, password)
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

// startAuthServer starts an authentication proxy that forwards requests to the app.
func startAuthServer(targetPort, authPort, password string) error {
	ip := network.GetLocalIP()
	if ip == "" {
		return fmt.Errorf("could not determine local IP address")
	}

	target, err := url.Parse("http://localhost:" + targetPort)
	if err != nil {
		return err
	}
	// ReverseProxy sets the upstream Host correctly and upgrades WebSocket
	// connections (Vite/Next HMR), which the old hand-rolled proxy dropped.
	proxy := httputil.NewSingleHostReverseProxy(target)
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = target.Host // present the upstream's own host (dev-server host checks)
	}

	var handler http.Handler = proxy
	if password != "" {
		handler = middleware.AuthMiddleware(proxy, password)
		fmt.Printf("🔐 Authentication enabled - Password required to access the app\n")
	}

	server := &http.Server{Addr: "0.0.0.0:" + authPort, Handler: handler}

	fmt.Printf("🔗 Auth Proxy: http://%s:%s\n", ip, authPort)
	qrcode.GenerateQrCodeWithMessage(ip+":"+authPort, "📱 Scan this on your phone:")

	return server.ListenAndServe()
}

// runWithInstallRetry runs the app, installing dependencies once and retrying if
// the first attempt can't start. QR code is shown only after the server is up.
func runWithInstallRetry(pm string, scripts, flags []string, port, password string) error {
	if password != "" {
		authPort := strconv.Itoa(getAvailablePort(port))
		go func() {
			if err := startAuthServer(port, authPort, password); err != nil {
				fmt.Printf("❌ Auth server error: %v\n", err)
			}
		}()
	}

	if started, err := tryScripts(pm, scripts, flags, port, password); started {
		return err
	}

	fmt.Println("Installing dependencies...")
	install := exec.Command(pm, "install")
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr
	if err := install.Run(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	if started, err := tryScripts(pm, scripts, flags, port, password); started {
		return err
	}
	return fmt.Errorf("failed to start app with %s", pm)
}

// tryScripts runs each candidate script until one actually binds the dev port.
// Returns started=true (with that run's exit error) once a server comes up, or
// started=false if every candidate exited before binding.
func tryScripts(pm string, scripts, flags []string, port, password string) (bool, error) {
	showQR := func() {
		if password == "" {
			qrcode.GenerateQrCodeWithMessage(network.GetLocalIP()+":"+port, "📱 Scan this on your phone:")
		}
	}
	for _, script := range scripts {
		cmd := exec.Command(pm, pmRunArgs(pm, script, flags)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			continue
		}
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()

		// ponytail: 15s bind window; a slower build past that is assumed to be the server.
		deadline := time.After(15 * time.Second)
		bound := false
		for !bound {
			select {
			case <-done:
				bound = true // exited before binding → try the next candidate
			case <-deadline:
				showQR()
				return true, <-done
			default:
				if portOpen(port) {
					showQR()
					return true, <-done
				}
				time.Sleep(200 * time.Millisecond)
			}
		}
	}
	return false, nil
}

// portOpen reports whether something is already listening on the local port.
func portOpen(port string) bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 200*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// pmRunArgs builds the args to run a package.json script with flags.
// npm needs "run <script>" and a "--" separator before script flags; yarn and
// pnpm forward extra args to the script directly.
func pmRunArgs(pm, script string, flags []string) []string {
	if pm == "npm" {
		args := []string{"run", script}
		if len(flags) > 0 {
			args = append(args, "--")
			args = append(args, flags...)
		}
		return args
	}
	return append([]string{script}, flags...)
}

// getAvailablePort finds an available port starting just above the given port.
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
			if err := listener.Close(); err != nil {
				return 0
			}
			return testPort
		}
	}
	return startSearchPort + 1
}

// runVite runs a Vite-based app (React or Vue) — identical bar the label.
func runVite(appName, pm, port, password string) error {
	fmt.Printf("🚀 Starting %s app...\n", appName)
	ip := network.GetLocalIP()
	if port == "" {
		port = "5173" // default Vite port
	}
	host := "0.0.0.0"
	if password != "" {
		host = "127.0.0.1"
	}
	fmt.Printf("Local:   http://localhost:%s\n", port)
	if password == "" {
		fmt.Printf("Network: http://%s:%s\n", ip, port)
	}
	// dev first (Vite), start second (CRA) as a fallback
	return runWithInstallRetry(pm, []string{"dev", "start"}, []string{"--port", port, "--host", host}, port, password)
}

func runNextJS(pm, port, password string) error {
	fmt.Println("🚀 Starting Next.js app...")
	ip := network.GetLocalIP()
	if port == "" {
		port = "3000" // default Next.js port
	}
	host := "0.0.0.0"
	if password != "" {
		host = "127.0.0.1"
	}
	fmt.Printf("Local:   http://localhost:%s\n", port)
	if password == "" {
		fmt.Printf("Network: http://%s:%s\n", ip, port)
	}
	return runWithInstallRetry(pm, []string{"dev"}, []string{"--port", port, "-H", host}, port, password)
}

func runGo(password string) error {
	fmt.Println("🚀 Starting Go app...")
	if password != "" {
		fmt.Println("\n\033[33mWARNING: --password is not supported for Go projects.\033[0m")
		fmt.Println("DevShare can't place an auth proxy in front of a Go app (it doesn't manage its port/binding). Running without authentication.")
	}
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Go app: %w", err)
	}
	return nil
}

// runLaravel runs the Laravel app.
func runLaravel(port, password string) error {
	fmt.Println("🚀 Starting Laravel app...")
	ip := network.GetLocalIP()
	if port == "" {
		port = "8000" // default Laravel port
	}
	host := "0.0.0.0"
	if password != "" {
		host = "127.0.0.1"
	}
	cmd := exec.Command("php", "artisan", "serve", "--host", host, "--port", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Local:   http://localhost:%s\n", port)
	if password != "" {
		authPort := strconv.Itoa(getAvailablePort(port))
		go func() {
			if err := startAuthServer(port, authPort, password); err != nil {
				fmt.Printf("❌ Auth server error: %v\n", err)
			}
		}()
	} else {
		fmt.Printf("Network: http://%s:%s\n", ip, port)
		qrcode.GenerateQrCodeWithMessage(ip+":"+port, "📱 Scan this on your phone:")
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Laravel app: %w", err)
	}
	return nil
}

// runNodeJS runs the Node.js app.
func runNodeJS(pm, port, password string) error {
	fmt.Println("🚀 Starting Node.js app...")
	ip := network.GetLocalIP()
	if port == "" {
		port = "3000" // default Node.js port
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
	}

	fmt.Println("\n\033[33mWARNING: Your Node.js app may be listening on all interfaces (0.0.0.0).\033[0m")
	fmt.Println("For security, ensure your app binds to 127.0.0.1 to prevent bypassing authentication.")
	fmt.Println("If you control the app, update your server code to listen only on 127.0.0.1, for example:")
	fmt.Println()
	fmt.Println("// Node.js (Express example):")
	fmt.Println("const host = process.env.HOST || '127.0.0.1';")
	fmt.Println("const port = process.env.PORT || 3000;")
	fmt.Println("app.listen(port, host, () => {")
	fmt.Println("  console.log(`Server running at http://${host}:${port}/`);")
	fmt.Println("});")

	tryAll := func() bool {
		for _, args := range pmCmds {
			fmt.Printf("Trying: %s %s\n", pm, args)
			printNetworkInfo()
			cmd := exec.Command(pm, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err == nil {
				return true
			} else {
				fmt.Printf("⚠️  Failed to run %s %s: %v\n", pm, args, err)
			}
		}
		for _, entry := range entryFiles {
			if _, err := os.Stat(entry.file); err != nil {
				continue
			}
			bin := "node"
			if entry.useTs {
				bin = "ts-node"
			}
			fmt.Printf("Trying: %s %s\n", bin, entry.file)
			printNetworkInfo()
			cmd := exec.Command(bin, entry.file)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err == nil {
				return true
			} else {
				fmt.Printf("⚠️  Failed to run %s %s: %v\n", bin, entry.file, err)
			}
		}
		return false
	}

	if tryAll() {
		return nil
	}
	fmt.Println("Installing dependencies...")
	install := exec.Command(pm, "install")
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr
	if err := install.Run(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}
	if tryAll() {
		return nil
	}
	return fmt.Errorf("could not start Node.js app: no working package manager script or entry file")
}
