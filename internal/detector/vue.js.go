package detector

import (
	"encoding/json"
	"os"
)

// isVueJSProject checks if the given directory is a Vue.js project
func isVueJSProject(dir string) (bool, string) {
	packagePath := dir + string(os.PathSeparator) + "package.json"

	if _, err := os.Stat(packagePath); err != nil {
		return false, ""
	}

	pkgManager, err := packageManagerDetector(dir)
	if err != nil {
		return false, ""
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

	// check vue in dependencies or devDependencies
	if pkg.Dependencies != nil {
		if _, ok := pkg.Dependencies["vue"]; ok {
			return true, pkgManager
		}
	}

	if pkg.DevDependencies != nil {
		if _, ok := pkg.DevDependencies["vue"]; ok {
			return true, pkgManager
		}
	}

	return false, ""
}
