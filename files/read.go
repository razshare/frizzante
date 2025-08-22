package files

import (
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
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()

	buf := make([]byte, max)

	var size int
	for {
		if size, err = file.Read(buf); err != nil {
			return
		}

		if size == 0 {
			return
		}

		if size < max {
			if err = call(buf[:size-1]); err != nil {
				return
			}
		}

		if err = call(buf); err != nil {
			return
		}
	}
}
