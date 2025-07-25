package detector

import (
	"encoding/json"
	"os"
)

// isReactProject checks if the given directory is a React project
func isReactProject(dir string) (bool, string) {
	packagePath := dir + string(os.PathSeparator) + "package.json"

	// Check if package.json exists
	if _, err := os.Stat(packagePath); err != nil {
		return false, ""
	}

	// Check witch package manager is used
	pkgManager, err := packageManagerDetector(dir)
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(packagePath)
	if err != nil {
		return false, ""
	}

	type pkgJson struct {
		Dependencies    map[string]interface{} `json:"dependencies"`
		DevDependencies map[string]interface{} `json:"devDependencies"`
	}

	var pkg pkgJson
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false, ""
	}

	// Check for 'react' in dependencies or devDependencies
	if pkg.Dependencies != nil {
		if _, ok := pkg.Dependencies["react"]; ok {
			return true, pkgManager
		}
	}

	if pkg.DevDependencies != nil {
		if _, ok := pkg.DevDependencies["react"]; ok {
			return true, pkgManager
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
			return true, pkgManager
		}
	}

	return false, ""
}
