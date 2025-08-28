package embeds

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"slices"
)

func ReadDirectory(efs embed.FS, from string) (entries []string, err error) {
	entries = make([]string, 0)
	var names []fs.DirEntry
	if names, err = efs.ReadDir(from); err != nil {
		return
	}

	for _, name := range names {
		if name.IsDir() {
			var subentries []string
			subentries, err = ReadDirectory(efs, fmt.Sprintf("%s/%s", from, name.Name()))
			if err != nil {
				return
			}

			entries = slices.Concat(entries, subentries)
			continue
		}

		entries = append(entries, fmt.Sprintf("%s/%s", from, name.Name()))
	}

	return
}

// ReadFileInChunks reads a file in chunks of a set maximum size.
func ReadFileInChunks(efs embed.FS, from string, max int, call func([]byte) error) (err error) {
	var file fs.File
	if file, err = efs.Open(from); err != nil {
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
