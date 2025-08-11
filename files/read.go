package files

import (
	"fmt"
	"os"
	"slices"
)

func ReadDirectory(dn string) ([]string, error) {
	items := make([]string, 0)
	entries, readDirError := os.ReadDir(dn)
	if readDirError != nil {
		return nil, readDirError
	}

	for _, entry := range entries {
		if entry.IsDir() {
			items2, readDirLocalError := ReadDirectory(fmt.Sprintf("%s/%s", dn, entry.Name()))
			if readDirLocalError != nil {
				return nil, readDirLocalError
			}

			items = slices.Concat(items, items2)
			continue
		}

		items = append(items, fmt.Sprintf("%s/%s", dn, entry.Name()))
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
