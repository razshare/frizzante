package embeds

import (
	"embed"
	"fmt"
	"io/fs"
	"slices"
)

func ReadDirectory(efs embed.FS, n string) ([]string, error) {
	its := make([]string, 0)
	ents, err := efs.ReadDir(n)
	if err != nil {
		return nil, err
	}

	for _, ent := range ents {
		if ent.IsDir() {
			var sits []string
			sits, err = ReadDirectory(efs, fmt.Sprintf("%s/%s", n, ent.Name()))
			if err != nil {
				return nil, err
			}

			its = slices.Concat(its, sits)
			continue
		}

		its = append(its, fmt.Sprintf("%s/%s", n, ent.Name()))
	}

	return its, nil
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
