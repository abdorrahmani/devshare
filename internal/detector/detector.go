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

func IsLaravelProject(dir string) bool {
	// Check for artisan file
	artisanPath := dir + string(os.PathSeparator) + "artisan"
	if _, err := os.Stat(artisanPath); err == nil {
		return true
	}

	// Check for composer.json and laravel/framework
	composerPath := dir + string(os.PathSeparator) + "composer.json"
	vendorLaravelPath := dir + string(os.PathSeparator) + "vendor" + string(os.PathSeparator) + "laravel" + string(os.PathSeparator) + "framework"
	if _, err := os.Stat(composerPath); err == nil {
		if _, err := os.Stat(vendorLaravelPath); err == nil {
			return true
		}
	}
	return false
}

// GetFilesInDir checks the directory and prints if it's a Laravel project.
func GetFilesInDir(dir string) {
	if dir == "" {
		fmt.Println("No directory provided.")
		return
	}

	fmt.Println("Scanning directory:", dir)
	if IsLaravelProject(dir) {
		fmt.Println("✅ Laravel project detected!")
		fmt.Println("You can now share your Laravel project over LAN.")
	} else {
		fmt.Println("❌ Not a Laravel project.")
	}
}
