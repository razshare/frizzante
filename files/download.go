package files

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile downloads a file to the disk.
func DownloadFile(url string, name string) (err error) {
	var res *http.Response
	if res, err = http.Get(url); err != nil {
		return
	}

	var data []byte
	if data, err = io.ReadAll(res.Body); err != nil {
		return
	}

	dir := filepath.Dir(name)
	if !IsDirectory(dir) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return
		}
	}

	if err = os.WriteFile(name, data, os.ModePerm); err != nil {
		return
	}

	return nil
}
