package cli

import (
	"errors"
	"os"
)

// isDirectory checks if file exists and is a directory.
func isDirectory(fileName string) bool {
	stat, statError := os.Stat(fileName)
	if statError != nil {
		return !errors.Is(statError, os.ErrNotExist)
	}
	return stat.IsDir()
}
