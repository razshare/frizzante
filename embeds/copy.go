package embeds

import (
	"embed"
	"github.com/razshare/frizzante/files"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyFile(efs embed.FS, from string, to string) error {
	var err error
	var src fs.File
	src, err = efs.Open(from)
	if err != nil {
		return err
	}

	todir := filepath.Dir(to)

	if !files.IsDirectory(todir) {
		err = os.MkdirAll(todir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if files.IsFile(to) {
		err = os.Remove(to)
		if err != nil {
			return err
		}
	}

	var dst *os.File
	dst, err = os.Create(to)
	if err != nil {
		_ = src.Close()
		return err
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		_ = src.Close()
		_ = dst.Close()
		return err
	}

	err = src.Close()
	if err != nil {
		return err
	}

	err = dst.Close()
	if err != nil {
		return err
	}

	return nil
}

func CopyDirectory(efs embed.FS, from string, to string) error {
	froms, err := ReadDirectory(efs, from)
	if err != nil {
		return err
	}

	for _, f := range froms {
		err = CopyFile(efs, f, filepath.Join(to, f))
		if err != nil {
			return err
		}
	}

	return nil
}
