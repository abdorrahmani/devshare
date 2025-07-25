package detector

import "os"

// isGoProject checks if the given directory is a Go project
func isGoProject(dir string) bool {
	// check for go.mod file
	goModPath := dir + string(os.PathSeparator) + "go.mod"
	if _, err := os.Stat(goModPath); err == nil {
		return true
	}

	// check for main.go file
	mainGoPath := dir + string(os.PathSeparator) + "main.go"
	if _, err := os.Stat(mainGoPath); err == nil {
		return true
	}
	return false
}
