package files

import (
	"errors"
	"os"
)

// DeleteFile deletes a file from the disk.
func DeleteFile(n string) bool {
	e := os.Remove(n)
	return nil == e || !errors.Is(e, os.ErrNotExist)
}
