package files

import "os"

// IsFile check if file exists and is a file.
func IsFile(name string) bool {
	if stat, err := os.Stat(name); err == nil {
		return !stat.IsDir()
	}
	return false
}
