package files

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile downloads a file to the disk.
func DownloadFile(url string, fileName string) (err error) {
	var response *http.Response
	if response, err = http.Get(url); err != nil {
		return
	}
	if response.Body != nil {
		defer func() {
			if cerr := response.Body.Close(); cerr != nil {
				if err == nil {
					err = cerr
				}
			}
		}()
	}
	var data []byte
	if data, err = io.ReadAll(response.Body); err != nil {
		return
	}
	directoryName := filepath.Dir(fileName)
	if !IsDirectory(directoryName) {
		if err = os.MkdirAll(directoryName, os.ModePerm); err != nil {
			return
		}
	}
	err = os.WriteFile(fileName, data, os.ModePerm)
	return
}
