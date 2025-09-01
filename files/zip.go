package files

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ZipFile zips a file to the disk.
func ZipFile(from string, to string) (err error) {
	if err = os.MkdirAll(filepath.Dir(to), os.ModePerm); err != nil {
		return
	}

	var zipFile *os.File
	if zipFile, err = os.Create(to); err != nil {
		return
	}

	defer func() {
		if zipFile == nil {
			return
		}
		if cerr := zipFile.Close(); cerr != nil {
			err = cerr
		}
	}()

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if zipWriter == nil {
			return
		}
		if cerr := zipWriter.Close(); cerr != nil {
			err = cerr
		}
	}()

	var writer io.Writer
	if writer, err = zipWriter.Create(filepath.Base(from)); err != nil {
		return
	}

	var file *os.File
	file, err = os.Open(from)
	if err != nil {
		return
	}

	if _, err = io.Copy(writer, file); err != nil {
		if cerr := file.Close(); cerr != nil {
			err = cerr
			return
		}
		return
	}

	return nil
}

// ZipDirectory zips a directory to the disk.
func ZipDirectory(from string, to string) (err error) {
	if err = os.MkdirAll(filepath.Dir(to), os.ModePerm); err != nil {
		return
	}

	var zipFile *os.File
	zipFile, err = os.Create(to)
	if err != nil {
		return
	}

	defer func() {
		if zipFile == nil {
			return
		}
		if cerr := zipFile.Close(); cerr != nil {
			err = cerr
		}
	}()

	zipFileWriter := zip.NewWriter(zipFile)

	defer func() {
		if zipFileWriter == nil {
			return
		}
		if cerr := zipFileWriter.Close(); cerr != nil {
			err = cerr
		}
	}()

	err = filepath.Walk(from, func(name string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(name)
		if err != nil {
			return err
		}

		writer, err := zipFileWriter.Create(strings.TrimPrefix(name, from+"/"))
		if err != nil {
			if cerr := file.Close(); cerr != nil {
				return cerr
			}
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})

	return
}
