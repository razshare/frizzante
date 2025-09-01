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
		return
	}

	return nil
}

// ZipDirectory zips a directory to the disk.
func ZipDirectory(from string, to string) (err error) {
	err = os.MkdirAll(filepath.Dir(to), os.ModePerm)
	if err != nil {
		return
	}

	var z *os.File

	z, err = os.Create(to)
	if err != nil {
		return
	}

	defer func(z *os.File) { err = z.Close() }(z)

	zw := zip.NewWriter(z)

	defer func(zw *zip.Writer) { err = zw.Close() }(zw)

	err = filepath.Walk(from, func(p string, i fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if i.IsDir() {
			return nil
		}

		f, err := os.Open(p)
		if err != nil {
			return err
		}

		zww, err := zw.Create(strings.TrimPrefix(p, from+"/"))
		if err != nil {
			return err
		}

		_, err = io.Copy(zww, f)
		if err != nil {
			return err
		}

		return nil
	})

	return
}
