package fs

import (
	"embed"
	"errors"
	"os"
)

// ExistsInEmbeddedFileSystem checks if file exists.
func ExistsInEmbeddedFileSystem(efs embed.FS, fileName string) bool {
	return isEmbeddedFile(efs, fileName) || IsEmbeddedDirectory(efs, fileName)
}

// isEmbeddedFile check if file exists and is a file.
func isEmbeddedFile(efs embed.FS, fileName string) bool {
	_, err := efs.ReadFile(fileName)
	if err != nil {
		return false
	}
	return true
}

// IsEmbeddedDirectory checks if file exists and is a directory.
func IsEmbeddedDirectory(efs embed.FS, fileName string) bool {
	_, err := efs.ReadDir(fileName)
	if err != nil {
		return false
	}
	return true
}

// FileExists checks if file exists.
func FileExists(fileName string) bool {
	_, statError := os.Stat(fileName)
	return nil == statError || !errors.Is(statError, os.ErrNotExist)
}

// IsFile check if file exists and is a file.
func IsFile(fileName string) bool {
	stat, statError := os.Stat(fileName)
	if statError != nil {
		return !errors.Is(statError, os.ErrNotExist)
	}
	return !stat.IsDir()
}

// IsDirectory checks if file exists and is a directory.
func IsDirectory(fileName string) bool {
	stat, statError := os.Stat(fileName)
	if statError != nil {
		return !errors.Is(statError, os.ErrNotExist)
	}
	return stat.IsDir()
}
