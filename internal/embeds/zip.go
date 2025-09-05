package embeds

import (
	"archive/zip"
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ZipFile zips a file to the disk.
func ZipFile(efs embed.FS, from string, to string) (err error) {
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

	var ioWriter io.Writer
	if ioWriter, err = zipWriter.Create(filepath.Base(from)); err != nil {
		return
	}

	var file fs.File
	if file, err = efs.Open(from); err != nil {
		return
	}

	if _, err = io.Copy(ioWriter, file); err != nil {
		return
	}

	return nil
}

// ZipDirectory zips a directory to the disk.
func ZipDirectory(efs embed.FS, from string, to string) (err error) {
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

	err = filepath.Walk(from, func(name string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := efs.Open(name)
		if err != nil {
			return err
		}

		ioWriter, err := zipWriter.Create(strings.TrimPrefix(name, from+"/"))
		if err != nil {
			return err
		}

		if _, err = io.Copy(ioWriter, file); err != nil {
			return err
		}

		return nil
	})

	return
}
