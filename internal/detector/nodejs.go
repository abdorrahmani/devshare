package detector

import "os"

// isNodeJSProject checks if the given directory is a Node.js project
func isNodeJSProject(dir string) (bool, string) {
	// Check for package.json file
	packageJsonPath := dir + string(os.PathSeparator) + "package.json"
	if _, err := os.Stat(packageJsonPath); err == nil {
		pkgManager, err := packageManagerDetector(dir)
		if err != nil {
			panic(err)
		}
		return true, pkgManager
	}
	return false, ""
}
