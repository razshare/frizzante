package files

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// UnzipFile unzips a file to a directory on the disk.
func UnzipFile(from string, to string) (err error) {
	var reader *zip.ReadCloser
	if reader, err = zip.OpenReader(from); err != nil {
		return
	}
	defer func() {
		if cerr := reader.Close(); cerr != nil {
			err = cerr
		}
	}()

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	for _, file := range reader.File {
		name := filepath.Join(to, file.Name)
		if file.FileInfo().IsDir() && !IsDirectory(name) {
			if err = os.MkdirAll(name, os.ModePerm); err != nil {
				return
			}
			continue
		}

		dir := filepath.Dir(name)

		if dir == "." {
			continue
		}

		if !IsDirectory(dir) {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return
			}
		}

		var zipFile *os.File
		if zipFile, err = os.OpenFile(name, flags, file.Mode()); err != nil {
			return
		}

		var zipReader io.ReadCloser
		if zipReader, err = file.Open(); err != nil {
			return
		}

		if _, err = io.Copy(zipFile, zipReader); err != nil {
			return
		}

		if err = zipFile.Close(); err != nil {
			return
		}

		if err = zipReader.Close(); err != nil {
			return
		}
	}

	return nil
}
