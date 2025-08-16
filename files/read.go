package files

import (
	"fmt"
	"os"
	"slices"
)

func ReadDirectory(n string) ([]string, error) {
	items := make([]string, 0)
	entries, err := os.ReadDir(n)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			var subitems []string
			subitems, err = ReadDirectory(fmt.Sprintf("%s/%s", n, entry.Name()))
			if err != nil {
				return nil, err
			}

			items = slices.Concat(items, subitems)
			continue
		}

		items = append(items, fmt.Sprintf("%s/%s", n, entry.Name()))
	}

	return items, nil
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
