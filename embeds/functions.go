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
func IsFile(self embed.FS, fname string) bool {
	file, err := self.Open(fname)
	if err != nil {
		return false
	}
	stat, statError := file.Stat()
	if statError != nil {
		return false
	}
	return !stat.IsDir()
}

func ReadDir(self embed.FS, dirname string) ([]string, error) {
	items := make([]string, 0)
	entries, readDirError := self.ReadDir(dirname)
	if readDirError != nil {
		return nil, readDirError
	}

	for _, entry := range entries {
		if entry.IsDir() {
			items2, readDirLocalError := ReadDir(self, fmt.Sprintf("%s/%s", dirname, entry.Name()))
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
func IsDirectory(self embed.FS, fname string) bool {
	file, err := self.Open(fname)
	if err != nil {
		return false
	}
	stat, statError := file.Stat()
	if statError != nil {
		return false
	}
	return stat.IsDir()
}

func FileReader(self embed.FS, fname string) (*bytes.Reader, os.FileInfo, error) {
	file, openError := self.Open(fname)
	if nil != openError {
		return nil, nil, openError
	}

	fileInfo, _ := file.Stat()

	buffer := make([]byte, fileInfo.Size())
	_, readError := file.Read(buffer)
	if nil != readError {
		closeError := file.Close()
		if nil != closeError {
			return nil, nil, closeError
		}
		return nil, nil, readError
	}

	closeError := file.Close()
	if nil != closeError {
		return nil, nil, closeError
	}
	return bytes.NewReader(buffer), fileInfo, nil
}

// ReadFileInChunks reads a file in chunks.
func ReadFileInChunks(self embed.FS, fname string, chunk int, fun func([]byte)) (err error) {
	file, openError := self.Open(fname)
	if nil != openError {
		return openError
	}
	defer func(file fs.File) { err = file.Close() }(file)

	buffer := make([]byte, chunk)

	for {
		count, readError := file.Read(buffer)
		if nil != readError {
			return readError
		}
		if count == 0 {
			return nil
		}
		if count < chunk {
			fun(buffer[:count-1])
		}
		fun(buffer)
	}
}
