package files

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
)

func ReadDirectory(dn string) ([]string, error) {
	its := make([]string, 0)
	ents, err := os.ReadDir(dn)
	if err != nil {
		return nil, err
	}

	for _, ent := range ents {
		n := dn + string(filepath.Separator) + ent.Name()
		if ent.IsDir() {
			var sits []string
			sits, err = ReadDirectory(n)
			if err != nil {
				return nil, err
			}

			its = slices.Concat(its, sits)
			continue
		}

		its = append(its, n)
	}

	return its, nil
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
