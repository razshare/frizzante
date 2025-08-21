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

// ReadFileInChunks reads a file in chunks.
func ReadFileInChunks(n string, c int, cb func([]byte) error) (err error) {
	file, openError := os.Open(n)
	if openError != nil {
		return openError
	}
	defer func(file *os.File) { err = file.Close() }(file)

	buffer := make([]byte, c)

	for {
		count, readError := file.Read(buffer)
		if readError != nil {
			return readError
		}
		if count == 0 {
			return nil
		}
		if count < c {
			callError := cb(buffer[:count-1])
			if callError != nil {
				return callError
			}
		}
		callError := cb(buffer)
		if callError != nil {
			return callError
		}
	}
}
