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

// DetectProjectType returns the project type (react, nextjs,vue.js, go, nodejs) and package manager if applicable.
func DetectProjectType(dir string) (string, string) {
	if isLaravelProject(dir) {
		fmt.Println("✅ Laravel project detected!")
		return "laravel", ""
	} else if ok, pkgManager := isNextJSProject(dir); ok {
		fmt.Println("✅ Next.js project detected!")
		fmt.Printf("📦 Using package manager: %s\n", pkgManager)
		return "nextjs", pkgManager
	} else if ok, pkgManager := isReactProject(dir); ok {
		fmt.Println("✅ React project detected!")
		fmt.Printf("📦 Using package manager: %s\n", pkgManager)
		return "react", pkgManager
	} else if ok, pkgManager := isVueJSProject(dir); ok {
		fmt.Println("✅ Vue.js project detected!")
		fmt.Printf("📦 Using package manager: %s\n", pkgManager)
		return "vue", pkgManager
	} else if isGoProject(dir) {
		fmt.Println("✅ Go project detected!")
		return "go", ""
	} else if ok, pkgManager := isNodeJSProject(dir); ok {
		fmt.Println("✅ Node.js project detected!")
		fmt.Printf("📦 Using package manager: %s\n", pkgManager)
		return "nodejs", pkgManager
	}
	return "", ""
}
