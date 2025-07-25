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

// GetFilesInDir checks the directory and prints if it's a Laravel/React project.
func GetFilesInDir(dir string) {
	if dir == "" {
		fmt.Println("No directory provided.")
		return
	}

	fmt.Println("Scanning directory:", dir)
	if isLaravelProject(dir) {
		fmt.Println("✅ Laravel project detected!")
		fmt.Println("You can now share your Laravel project over LAN.")
	} else if isNextJSProject(dir) {
		fmt.Println("✅ Next.js project detected!")
		fmt.Println("You can now share your Next.js project over LAN.")
	} else if isReactProject(dir) {
		fmt.Println("✅ React project detected!")
		fmt.Println("You can now share your React project over LAN.")
	} else if isGoProject(dir) {
		fmt.Println("✅ Go project detected!")
		fmt.Println("You can now share your Go project over LAN.")
	} else if isNodeJSProject(dir) {
		fmt.Println("✅ Node.js project detected!")
		fmt.Println("You can now share your Node.js project over LAN.")
	} else {
		fmt.Println("❌ Not a Laravel/Go/React project.")
	}
}
