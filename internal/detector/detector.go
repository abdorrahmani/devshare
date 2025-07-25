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

// GetFilesInDir checks the directory and prints if it's a Laravel/Go/Next.js/React/Node.js project.
func GetFilesInDir(dir string) string {
	if dir == "" {
		fmt.Println("No directory provided.")
		return ""
	}

	fmt.Println("Scanning directory:", dir)
	if isLaravelProject(dir) {
		fmt.Println("✅ Laravel project detected!")
		fmt.Println("You can now share your Laravel project over LAN.")
	} else if ok, pkgManager := isNextJSProject(dir); ok == true {
		fmt.Println("✅ Next.js project detected!")
		fmt.Println("You can now share your Next.js project over LAN.")
		fmt.Printf("🚀 Using package manager: %s\n", pkgManager)
		return pkgManager
	} else if ok, pkgManager := isReactProject(dir); ok == true {
		fmt.Println("✅ React project detected!")
		fmt.Println("You can now share your React project over LAN.")
		fmt.Printf("🚀 Using package manager: %s\n", pkgManager)
		return pkgManager
	} else if isGoProject(dir) {
		fmt.Println("✅ Go project detected!")
		fmt.Println("You can now share your Go project over LAN.")
	} else if ok, pkgManager := isNodeJSProject(dir); ok == true {
		fmt.Println("✅ Node.js project detected!")
		fmt.Println("You can now share your Node.js project over LAN.")
		fmt.Printf("🚀 Using package manager: %s\n", pkgManager)
		return pkgManager
	} else {
		fmt.Println("❌ Not a Laravel/Go/Node.js/Next.js/React project.")
	}

	return ""
}
