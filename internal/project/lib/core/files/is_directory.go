package files

import "os"

// IsDirectory checks if file exists and is a directory.
func IsDirectory(name string) bool {
	if stat, err := os.Stat(name); err == nil {
		return stat.IsDir()
	}
	return false
}
