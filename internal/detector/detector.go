package detector

import (
	"encoding/json"
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

func isReactProject(dir string) bool {
	packagePath := dir + string(os.PathSeparator) + "package.json"

	// Check if package.json exists
	if _, err := os.Stat(packagePath); err != nil {
		return false
	}

	// Read and parse package.json
	data, err := os.ReadFile(packagePath)
	if err != nil {
		return false
	}

	type pkgJson struct {
		Dependencies    map[string]interface{} `json:"dependencies"`
		DevDependencies map[string]interface{} `json:"devDependencies"`
	}
	var pkg pkgJson
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}

	// Check for 'react' in dependencies or devDependencies
	if pkg.Dependencies != nil {
		if _, ok := pkg.Dependencies["react"]; ok {
			return true
		}
	}
	if pkg.DevDependencies != nil {
		if _, ok := pkg.DevDependencies["react"]; ok {
			return true
		}
	}

	// Check for typical React entry files
	entryFiles := []string{
		"src/App.js", "src/App.jsx", "src/App.tsx",
		"src/index.js", "src/index.jsx", "src/index.tsx",
	}
	for _, f := range entryFiles {
		entryPath := dir + string(os.PathSeparator) + f
		if _, err := os.Stat(entryPath); err == nil {
			return true
		}
	}

	return false
}

// GetFilesInDir checks the directory and prints if it's a Laravel/Reacy project.
func GetFilesInDir(dir string) {
	if dir == "" {
		fmt.Println("No directory provided.")
		return
	}

	fmt.Println("Scanning directory:", dir)
	if IsLaravelProject(dir) {
		fmt.Println("✅ Laravel project detected!")
		fmt.Println("You can now share your Laravel project over LAN.")
	} else if isReactProject(dir) {
		fmt.Println("✅ React project detected!")
		fmt.Println("You can now share your React project over LAN.")
	} else {
		fmt.Println("❌ Not a Laravel/React project.")
	}
}
