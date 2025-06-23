package libfs

import (
	"archive/zip"
	"bytes"
	"embed"
	"errors"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// EfsIsFile check if file exists and is a file.
func EfsIsFile(efs embed.FS, fileName string) bool {
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

// EfsIsDirectory checks if file exists and is a directory.
func EfsIsDirectory(efs embed.FS, fileName string) bool {
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

// IsFile check if file exists and is a file.
func IsFile(fileName string) bool {
	stat, statError := os.Stat(fileName)
	if statError != nil {
		return false
	}
	return !stat.IsDir()
}

// IsDirectory checks if file exists and is a directory.
func IsDirectory(fileName string) bool {
	stat, statError := os.Stat(fileName)
	if statError != nil {
		return false
	}
	return stat.IsDir()
}

// DeleteFile deletes a file from the disk.
func DeleteFile(fileName string) bool {
	removeError := os.Remove(fileName)
	return nil == removeError || !errors.Is(removeError, os.ErrNotExist)
}

// UnzipFile unzips a file to the disk.
func UnzipFile(zipFileName string, directoryName string) (err error) {
	zipReader, zipOpenError := zip.OpenReader(zipFileName)
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

// ZipFile zips a file to the disk.
func ZipFile(fileName string, zipFileName string) (err error) {
	mkdirError := os.MkdirAll(filepath.Dir(zipFileName), os.ModePerm)
	if mkdirError != nil {
		return mkdirError
	}

	archive, createError := os.Create(zipFileName)
	if createError != nil {
		return createError
	}
	defer func(archive *os.File) { err = archive.Close() }(archive)

	writer := zip.NewWriter(archive)
	defer func(writer *zip.Writer) { err = writer.Close() }(writer)

	entry, entryError := writer.Create(filepath.Base(fileName))
	if entryError != nil {
		return entryError
	}

	file, fileError := os.Open(fileName)
	if fileError != nil {
		return fileError
	}

	_, copyError := io.Copy(entry, file)
	if copyError != nil {
		return copyError
	}

	return nil
}

// ZipDirectory zips a directory to the disk.
func ZipDirectory(directoryName string, zipFileName string) (err error) {
	mkdirError := os.MkdirAll(filepath.Dir(zipFileName), os.ModePerm)
	if mkdirError != nil {
		return mkdirError
	}

	archive, archiveError := os.Create(zipFileName)
	if archiveError != nil {
		return archiveError
	}
	defer func(archive *os.File) { err = archive.Close() }(archive)

	writer := zip.NewWriter(archive)
	defer func(writer *zip.Writer) { err = writer.Close() }(writer)

	return filepath.Walk(directoryName, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, openError := os.Open(path)
		if openError != nil {
			return openError
		}

		entry, entryError := writer.Create(strings.TrimPrefix(path, directoryName+"/"))
		if entryError != nil {
			return entryError
		}

		_, copyError := io.Copy(entry, file)
		if copyError != nil {
			return copyError
		}

		return nil
	})
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

	parentName := filepath.Dir(fileName)
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

func EfsFileReader(efs embed.FS, fileName string) (*bytes.Reader, os.FileInfo, error) {
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
	return bytes.NewReader(buffer), fileInfo, nil
}

func FileReader(fileName string) (*bytes.Reader, os.FileInfo, error) {
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

	return bytes.NewReader(buffer), fileInfo, nil
}

// ReadFileInChunks reads a file in chunks.
func ReadFileInChunks(fileName string, chunkSize int, callback func(data []byte) error) (err error) {
	file, openError := os.Open(fileName)
	if nil != openError {
		return openError
	}
	defer func(file *os.File) { err = file.Close() }(file)

	buffer := make([]byte, chunkSize)

	for {
		count, readError := file.Read(buffer)
		if nil != readError {
			return readError
		}
		if count == 0 {
			return nil
		}
		if count < chunkSize {
			callError := callback(buffer[:count-1])
			if callError != nil {
				return callError
			}
		}
		callError := callback(buffer)
		if callError != nil {
			return callError
		}
	}
}

// EfsReadFileInChunks reads a file in chunks.
func EfsReadFileInChunks(efs embed.FS, fileName string, chunkSize int, callback func([]byte)) (err error) {
	file, openError := efs.Open(fileName)
	if nil != openError {
		return openError
	}
	defer func(file fs.File) { err = file.Close() }(file)

	buffer := make([]byte, chunkSize)

	for {
		count, readError := file.Read(buffer)
		if nil != readError {
			return readError
		}
		if count == 0 {
			return nil
		}
		if count < chunkSize {
			callback(buffer[:count-1])
		}
		callback(buffer)
	}
}
