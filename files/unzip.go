package files

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// UnzipFile unzips a file to a directory on the disk.
func UnzipFile(zn string, dn string) (err error) {
	var r *zip.ReadCloser
	r, err = zip.OpenReader(zn)
	if err != nil {
		return
	}

	defer func(r *zip.ReadCloser) { err = r.Close() }(r)

	for _, zf := range r.File {
		zfn := filepath.Join(dn, zf.Name)
		if zf.FileInfo().IsDir() && !IsDirectory(zfn) {
			err = os.MkdirAll(zfn, os.ModePerm)
			if err != nil {
				return
			}
			continue
		}

		d := filepath.Dir(zfn)

		if d == "." {
			continue
		}

		if !IsDirectory(d) {
			err = os.MkdirAll(d, os.ModePerm)
			if err != nil {
				return
			}
		}

		var f *os.File

		f, err = os.OpenFile(zfn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode())
		if err != nil {
			return
		}

		var zfr io.ReadCloser

		zfr, err = zf.Open()
		if err != nil {
			return
		}

		if _, err = io.Copy(f, zfr); err != nil {
			return
		}

		_ = f.Close()
		_ = zfr.Close()
	}

	return nil
}
