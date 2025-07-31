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
func IsFile(efs embed.FS, fileName string) bool {
	file, err := efs.Open(fileName)
	if err != nil {
		return false
	}
	stat, statError := file.Stat()
	if statError != nil {
		return false
	}
	return !stat.IsDir()
}

func ReadDir(efs embed.FS, dirname string) ([]string, error) {
	items := make([]string, 0)
	entries, readDirError := efs.ReadDir(dirname)
	if readDirError != nil {
		return nil, readDirError
	}

	for _, entry := range entries {
		if entry.IsDir() {
			items2, readDirLocalError := ReadDir(efs, fmt.Sprintf("%s/%s", dirname, entry.Name()))
			if readDirLocalError != nil {
				return nil, readDirLocalError
			}

			items = slices.Concat(items, items2)
			continue
		}

		items = append(items, fmt.Sprintf("%s/%s", dirname, entry.Name()))
	}

	return items, nil
}

// IsDirectory checks if file exists and is a directory.
func IsDirectory(efs embed.FS, fileName string) bool {
	file, err := efs.Open(fileName)
	if err != nil {
		return false
	}
	stat, statError := file.Stat()
	if statError != nil {
		return false
	}
	return stat.IsDir()
}

func FileReader(efs embed.FS, fileName string) (*bytes.Reader, os.FileInfo, error) {
	file, openError := efs.Open(fileName)
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
func ReadFileInChunks(efs embed.FS, fileName string, chunk int, callback func([]byte)) (err error) {
	file, openError := efs.Open(fileName)
	if openError != nil {
		return openError
	}
	defer func(file fs.File) { err = file.Close() }(file)

	buffer := make([]byte, chunk)

	for {
		count, readError := file.Read(buffer)
		if readError != nil {
			return readError
		}
		if count == 0 {
			return nil
		}
		if count < chunk {
			callback(buffer[:count-1])
		}
		callback(buffer)
	}
}
