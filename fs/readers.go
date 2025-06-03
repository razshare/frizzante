package fs

import (
	"bytes"
	"embed"
	"os"
)

func ReaderFromEmbeddedFileName(efs embed.FS, fileName string) (*bytes.Reader, *os.FileInfo, error) {
	file, openError := efs.Open(fileName)
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
	return bytes.NewReader(buffer), &fileInfo, nil
}

func ReaderFromFileName(fileName string) (*bytes.Reader, *os.FileInfo, error) {
	file, openError := os.Open(fileName)
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
	return bytes.NewReader(buffer), &fileInfo, nil
}
