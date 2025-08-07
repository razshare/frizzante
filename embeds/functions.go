package embeds

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"slices"
)

// IsFile check if file exists and is a file.
func IsFile(efs embed.FS, n string) bool {
	file, err := efs.Open(n)
	if err != nil {
		return false
	}
	stat, statError := file.Stat()
	if statError != nil {
		return false
	}
	return !stat.IsDir()
}

// IsDirectory checks if file exists and is a directory.
func IsDirectory(efs embed.FS, n string) bool {
	file, err := efs.Open(n)
	if err != nil {
		return false
	}
	stat, statError := file.Stat()
	if statError != nil {
		return false
	}
	return stat.IsDir()
}

func ReadDir(efs embed.FS, dn string) ([]string, error) {
	items := make([]string, 0)
	entries, readDirError := efs.ReadDir(dn)
	if readDirError != nil {
		return nil, readDirError
	}

	for _, entry := range entries {
		if entry.IsDir() {
			items2, readDirLocalError := ReadDir(efs, fmt.Sprintf("%s/%s", dn, entry.Name()))
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

func FileReader(efs embed.FS, n string) (*bytes.Reader, os.FileInfo, error) {
	file, openError := efs.Open(n)
	if openError != nil {
		return nil, nil, openError
	}

	fileInfo, _ := file.Stat()

	buffer := make([]byte, fileInfo.Size())
	_, readError := file.Read(buffer)
	if readError != nil {
		closeError := file.Close()
		if closeError != nil {
			return nil, nil, closeError
		}
		return nil, nil, readError
	}

	closeError := file.Close()
	if closeError != nil {
		return nil, nil, closeError
	}
	return bytes.NewReader(buffer), fileInfo, nil
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
