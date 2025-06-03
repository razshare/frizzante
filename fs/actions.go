package fs

import (
	"errors"
	"os"
)

// DeleteFile deletes a file.
func DeleteFile(fileName string) bool {
	removeError := os.Remove(fileName)
	return nil == removeError || !errors.Is(removeError, os.ErrNotExist)
}
