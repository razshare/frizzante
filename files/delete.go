package files

import (
	"errors"
	"os"
)

// DeleteFile deletes a file from the disk.
func DeleteFile(name string) bool {
	err := os.Remove(name)
	return nil == err || !errors.Is(err, os.ErrNotExist)
}
