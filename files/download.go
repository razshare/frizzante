package files

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile downloads a file to the disk.
func DownloadFile(url string, n string) error {
	response, getError := http.Get(url)
	if getError != nil {
		return getError
	}

	data, readError := io.ReadAll(response.Body)
	if readError != nil {
		return readError
	}

	parentName := filepath.Dir(n)
	if !IsDirectory(parentName) {
		mkdirError := os.MkdirAll(parentName, os.ModePerm)
		if mkdirError != nil {
			return mkdirError
		}
	}

	writeErr := os.WriteFile(n, data, os.ModePerm)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
