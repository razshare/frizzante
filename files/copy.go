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
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return
		}
	}

	if IsFile(to) {
		if err = os.Remove(to); err != nil {
			return
		}
	}

	var fromFile *os.File
	if fromFile, err = os.Open(from); err != nil {
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
	if toFile, err = os.Create(to); err != nil {
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

	if _, err = io.Copy(toFile, fromFile); err != nil {
		return
	}

	return nil
}

func CopyDirectory(from string, to string) (err error) {
	var entries []string
	if entries, err = ReadDirectory(from); err != nil {
		return
	}

	for _, entry := range entries {
		name := filepath.Join(to, strings.TrimPrefix(entry, from))
		if err = CopyFile(entry, name); err != nil {
			return
		}
	}

	return
}
