package fs

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// DeleteFile deletes a file.
func DeleteFile(fileName string) bool {
	removeError := os.Remove(fileName)
	return nil == removeError || !errors.Is(removeError, os.ErrNotExist)
}

// UnzipFile unzips a file to a directory.
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

	_ = os.Remove(fileName)

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
