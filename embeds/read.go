package embeds

import (
	"embed"
	"fmt"
	"io/fs"
	"slices"
)

func ReadDirectory(efs embed.FS, dn string) ([]string, error) {
	items := make([]string, 0)
	entries, readDirError := efs.ReadDir(dn)
	if readDirError != nil {
		return nil, readDirError
	}

	for _, entry := range entries {
		if entry.IsDir() {
			items2, readDirLocalError := ReadDirectory(efs, fmt.Sprintf("%s/%s", dn, entry.Name()))
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
func ReadFileInChunks(efs embed.FS, n string, c int, cb func([]byte)) (err error) {
	file, openError := efs.Open(n)
	if openError != nil {
		return openError
	}
	defer func(file fs.File) { err = file.Close() }(file)

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
			cb(buffer[:count-1])
		}
		cb(buffer)
	}
}
