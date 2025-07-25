package detector

import "os"

// isNextJSProject checks if the given directory is a Next.js project
func isNextJSProject(dir string) bool {
	// Check for next.config.ts or js file
	nextConfigFiles := []string{
		"next.config.ts",
		"next.config.js",
	}

	for _, file := range nextConfigFiles {
		nextConfigPath := dir + string(os.PathSeparator) + file
		if _, err := os.Stat(nextConfigPath); err == nil {
			return true
		}
	}

	return false
}
