package detector

import (
	"fmt"
	"os"
)

func GetWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return ""
	}

	fmt.Println("Current working directory:", dir)
	return dir
}

// DetectProjectType returns the project type (react, nextjs, go, nodejs) and package manager if applicable.
func DetectProjectType(dir string) (string, string) {
	if isLaravelProject(dir) {
		fmt.Println("✅ Laravel project detected!")
		fmt.Println("You can now share your Laravel project over LAN.")
		return "laravel", ""
	} else if ok, pkgManager := isNextJSProject(dir); ok {
		fmt.Println("✅ Next.js project detected!")
		fmt.Println("You can now share your Next.js project over LAN.")
		fmt.Printf("📦 Using package manager: %s\n", pkgManager)
		return "nextjs", pkgManager
	} else if ok, pkgManager := isReactProject(dir); ok {
		fmt.Println("✅ React project detected!")
		fmt.Println("You can now share your React project over LAN.")
		fmt.Printf("📦 Using package manager: %s\n", pkgManager)
		return "react", pkgManager
	} else if isGoProject(dir) {
		fmt.Println("✅ Go project detected!")
		fmt.Println("You can now share your Go project over LAN.")
		return "go", ""
	} else if ok, pkgManager := isNodeJSProject(dir); ok {
		fmt.Println("✅ Node.js project detected!")
		fmt.Println("You can now share your Node.js project over LAN.")
		fmt.Printf("📦 Using package manager: %s\n", pkgManager)
		return "nodejs", pkgManager
	}
	return "", ""
}
