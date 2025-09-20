package embeds

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func CopyFile(efs embed.FS, from string, to string) (err error) {
	var src fs.File
	if src, err = efs.Open(from); err != nil {
		return
	}

	dir := filepath.Dir(to)

	if !files.IsDirectory(dir) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return
		}
	}

	if files.IsFile(to) {
		if err = os.Remove(to); err != nil {
			return
		}
	}

	var dst *os.File
	if dst, err = os.Create(to); err != nil {
		_ = src.Close()
		return
	}

	if _, err = io.Copy(dst, src); err != nil {
		_ = src.Close()
		_ = dst.Close()
		return
	}

	if err = src.Close(); err != nil {
		return
	}

	if err = dst.Close(); err != nil {
		return
	}

	return nil
}

func CopyDirectory(efs embed.FS, from string, to string) (err error) {
	var entries []string
	entries, err = ReadDirectory(efs, from)
	if err != nil {
		return
	}

	for _, ent := range entries {
		n := filepath.Join(to, strings.TrimPrefix(ent, from))
		err = CopyFile(efs, ent, n)
		if err != nil {
			return
		}
	}

	return nil
}
