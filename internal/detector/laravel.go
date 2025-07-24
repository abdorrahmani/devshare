package detector

import "os"

// IsLaravelProject checks if the given directory is a Laravel project
func isLaravelProject(dir string) bool {
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
