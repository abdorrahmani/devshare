package detector

import (
	"fmt"
	"os"
)

// packageManagerDetector  checks which package manager is used in the given directory
// by looking for lock files (yarn.lock, package-lock.json, or pnpm-lock.yaml).
func packageManagerDetector(dir string) (string, error) {
	packageManagers := map[string]string{
		"npm":  "package-lock.json",
		"yarn": "yarn.lock",
		"pnpm": "pnpm-lock.yaml",
	}

	for pkgManager, lockFile := range packageManagers {
		lockFilePath := dir + string(os.PathSeparator) + lockFile
		if _, err := os.Stat(lockFilePath); err == nil {
			return pkgManager, nil
		}
	}
	return "", fmt.Errorf("package manager not found in %s", dir)
}
