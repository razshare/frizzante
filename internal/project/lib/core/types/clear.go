package types

import (
	"os"
	"path/filepath"
)

// Clear removes the app/lib/types/server directory.
func Clear() (err error) {
	typesDirectoryName := filepath.Join("app", "lib", "types", "server")
	err = os.RemoveAll(typesDirectoryName)
	return
}
