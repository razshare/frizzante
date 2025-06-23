package efs

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
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
