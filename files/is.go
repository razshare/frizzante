package files

import "os"

// IsFile check if file exists and is a file.
func IsFile(n string) bool {
	stat, statError := os.Stat(n)
	if statError != nil {
		return false
	}
	return !stat.IsDir()
}

// IsDirectory checks if file exists and is a directory.
func IsDirectory(n string) bool {
	stat, statError := os.Stat(n)
	if statError != nil {
		return false
	}
	return stat.IsDir()
}
