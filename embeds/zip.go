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
func ZipFile(efs embed.FS, n string, zn string) (err error) {
	err = os.MkdirAll(filepath.Dir(zn), os.ModePerm)
	if err != nil {
		return
	}

	var z *os.File

	z, err = os.Create(zn)
	if err != nil {
		return
	}

	defer func(z *os.File) { err = z.Close() }(z)

	zw := zip.NewWriter(z)

	defer func(zw *zip.Writer) { err = zw.Close() }(zw)

	var zww io.Writer

	zww, err = zw.Create(filepath.Base(n))
	if err != nil {
		return
	}

	var f fs.File

	f, err = efs.Open(n)
	if err != nil {
		return
	}

	_, err = io.Copy(zww, f)
	if err != nil {
		return
	}

	return nil
}

// ZipDirectory zips a directory to the disk.
func ZipDirectory(efs embed.FS, dn string, zn string) (err error) {
	err = os.MkdirAll(filepath.Dir(zn), os.ModePerm)
	if err != nil {
		return
	}

	var z *os.File

	z, err = os.Create(zn)
	if err != nil {
		return
	}

	defer func(z *os.File) { err = z.Close() }(z)

	zw := zip.NewWriter(z)

	defer func(zw *zip.Writer) { err = zw.Close() }(zw)

	err = filepath.Walk(dn, func(p string, i fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if i.IsDir() {
			return nil
		}

		f, err := efs.Open(p)
		if err != nil {
			return err
		}

		zww, err := zw.Create(strings.TrimPrefix(p, dn+"/"))
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
