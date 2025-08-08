package files

import (
	"bytes"
	"os"
)

func NewFileReader(n string) (*bytes.Reader, os.FileInfo, error) {
	file, openError := os.Open(n)
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
