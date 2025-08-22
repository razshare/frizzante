package files

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyFile(from string, to string) (err error) {
	dir := filepath.Dir(to)

	if !IsDirectory(dir) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return
		}
	}

	if IsFile(to) {
		err = os.Remove(to)
		if err != nil {
			return
		}
	}

	var fromFile *os.File
	fromFile, err = os.Open(from)
	if err != nil {
		return
	}
	defer func() {
		if fromFile == nil {
			return
		}
		if cerr := fromFile.Close(); cerr != nil {
			err = cerr
		}
	}()

	var toFile *os.File
	toFile, err = os.Create(to)
	if err != nil {
		return
	}
	defer func() {
		if toFile == nil {
			return
		}
		if cerr := toFile.Close(); cerr != nil {
			err = cerr
		}
	}()

	_, err = io.Copy(toFile, fromFile)
	if err != nil {
		return
	}

	return nil
}

func CopyDirectory(from string, to string) (err error) {
	var ents []string
	if ents, err = ReadDirectory(from); err != nil {
		return
	}

	for _, ent := range ents {
		n := filepath.Join(to, strings.TrimPrefix(ent, from))
		err = CopyFile(ent, n)
		if err != nil {
			return
		}
	}

	return
}
