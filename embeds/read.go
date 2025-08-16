package embeds

import (
	"embed"
	"fmt"
	"io/fs"
	"slices"
)

func ReadDirectory(efs embed.FS, n string) ([]string, error) {
	items := make([]string, 0)
	entries, err := efs.ReadDir(n)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			var subitems []string
			subitems, err = ReadDirectory(efs, fmt.Sprintf("%s/%s", n, entry.Name()))
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
