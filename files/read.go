package files

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

func ReadDirectory(from string) (names []string, err error) {
	var entries []fs.DirEntry
	if entries, err = os.ReadDir(from); err != nil {
		return
	}

	for _, entry := range entries {
		name := from + string(filepath.Separator) + entry.Name()
		if entry.IsDir() {
			var subnames []string
			if subnames, err = ReadDirectory(name); err != nil {
				return
			}

			names = slices.Concat(names, subnames)

			continue
		}

		names = append(names, name)
	}

	return
}

// ReadFileInChunks reads a file in chunks of a set maximum size.
func ReadFileInChunks(name string, max int, call func([]byte) error) (err error) {
	var file *os.File
	if file, err = os.Open(name); err != nil {
		return
	}
	defer func() {
		if file == nil {
			return
		}
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()

	buf := make([]byte, max)

	var size int
	for {
		size, err = file.Read(buf)

		if size > 0 {
			if err = call(buf[:size]); err != nil {
				return
			}
		}

		if errors.Is(err, io.EOF) {
			err = nil
			return
		}
	}
}
