package fs

import (
	"archive/zip"
	"bytes"
	"embed"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// EfsFileExists checks if file (or directory) exists.
func EfsFileExists(efs embed.FS, fileName string) bool {
	return EfsIsFile(efs, fileName) || EfsIsDirectory(efs, fileName)
}

// EfsIsFile check if file exists and is a file.
func EfsIsFile(efs embed.FS, fileName string) bool {
	_, err := efs.ReadFile(fileName)
	if err != nil {
		return false
	}
	return true
}

// EfsIsDirectory checks if file exists and is a directory.
func EfsIsDirectory(efs embed.FS, fileName string) bool {
	_, err := efs.ReadDir(fileName)
	if err != nil {
		return false
	}
	return true
}

// FileExists checks if file (or directory) exists.
func FileExists(fileName string) bool {
	_, statError := os.Stat(fileName)
	return nil == statError || !errors.Is(statError, os.ErrNotExist)
}

// IsFile check if file exists and is a file.
func IsFile(fileName string) bool {
	stat, statError := os.Stat(fileName)
	if statError != nil {
		return !errors.Is(statError, os.ErrNotExist)
	}
	return !stat.IsDir()
}

// IsDirectory checks if file exists and is a directory.
func IsDirectory(fileName string) bool {
	stat, statError := os.Stat(fileName)
	if statError != nil {
		return !errors.Is(statError, os.ErrNotExist)
	}
	return stat.IsDir()
}

// DeleteFile deletes a file from the disk.
func DeleteFile(fileName string) bool {
	removeError := os.Remove(fileName)
	return nil == removeError || !errors.Is(removeError, os.ErrNotExist)
}

// UnzipFile unzips a file the disk.
func UnzipFile(fileName string, directoryName string) (err error) {
	zipReader, zipOpenError := zip.OpenReader(fileName)
	if zipOpenError != nil {
		log.Fatal(zipOpenError)
	}

	defer func(zipReader *zip.ReadCloser) {
		closeError := zipReader.Close()
		if closeError != nil {
			err = closeError
		}
	}(zipReader)

	for _, file := range zipReader.File {
		fileNameLocal := filepath.Join(directoryName, file.Name)
		fileIsDirectory := file.FileInfo().IsDir()
		fileIsDirectoryOnDisk := IsDirectory(fileNameLocal)

		if fileIsDirectory && !fileIsDirectoryOnDisk {
			mkdirError := os.MkdirAll(fileNameLocal, os.ModePerm)
			if mkdirError != nil {
				return mkdirError
			}
			continue
		}

		directoryNameLocal := filepath.Dir(fileNameLocal)

		if "." == directoryNameLocal {
			continue
		}

		parentExists := IsDirectory(directoryNameLocal)

		if !parentExists {
			mkdirError := os.MkdirAll(directoryNameLocal, os.ModePerm)
			if mkdirError != nil {
				return mkdirError
			}
		}

		dstFile, openError := os.OpenFile(fileNameLocal, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if openError != nil {
			return openError
		}

		fileInArchive, openError := file.Open()
		if openError != nil {
			return openError
		}

		if _, copyError := io.Copy(dstFile, fileInArchive); copyError != nil {
			return copyError
		}

		_ = dstFile.Close()
		_ = fileInArchive.Close()
	}

	return nil
}

// DownloadFile downloads a file to the disk.
func DownloadFile(url string, fileName string) error {
	response, getError := http.Get(url)
	if getError != nil {
		return getError
	}

	data, readError := io.ReadAll(response.Body)
	if readError != nil {
		return readError
	}

	parentName := filepath.Base(fileName)
	if !IsDirectory(parentName) {
		mkdirError := os.MkdirAll(parentName, os.ModePerm)
		if mkdirError != nil {
			return mkdirError
		}
	}

	writeErr := os.WriteFile(fileName, data, os.ModePerm)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

func EfsFileReader(efs embed.FS, fileName string) (*bytes.Reader, *os.FileInfo, error) {
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

func FileReader(fileName string) (*bytes.Reader, *os.FileInfo, error) {
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
