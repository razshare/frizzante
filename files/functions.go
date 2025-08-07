package files

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// IsFile check if file exists and is a file.
func IsFile(n string) bool {
	stat, statError := os.Stat(n)
	if statError != nil {
		return false
	}
	return !stat.IsDir()
}

// IsDirectory checks if file exists and is a directory.
func IsDirectory(n string) bool {
	stat, statError := os.Stat(n)
	if statError != nil {
		return false
	}
	return stat.IsDir()
}

// DeleteFile deletes a file from the disk.
func DeleteFile(n string) bool {
	e := os.Remove(n)
	return nil == e || !errors.Is(e, os.ErrNotExist)
}

// UnzipFile unzips a file to a directory on the disk.
func UnzipFile(zn string, dn string) (err error) {
	reader, readerError := zip.OpenReader(zn)
	if readerError != nil {
		log.Fatal(readerError)
	}

	defer func(reader *zip.ReadCloser) {
		closeError := reader.Close()
		if closeError != nil {
			err = closeError
		}
	}(reader)

	for _, f := range reader.File {
		fileNameLocal := filepath.Join(dn, f.Name)
		isDir := f.FileInfo().IsDir()
		isDirOnDisk := IsDirectory(fileNameLocal)

		if isDir && !isDirOnDisk {
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

		fileInDestination, openError := os.OpenFile(fileNameLocal, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if openError != nil {
			return openError
		}

		fileInArchive, openError := f.Open()
		if openError != nil {
			return openError
		}

		if _, copyError := io.Copy(fileInDestination, fileInArchive); copyError != nil {
			return copyError
		}

		_ = fileInDestination.Close()
		_ = fileInArchive.Close()
	}

	return nil
}

// ZipFile zips a file to a directory on the disk.
func ZipFile(n string, zn string) (err error) {
	mkdirError := os.MkdirAll(filepath.Dir(zn), os.ModePerm)
	if mkdirError != nil {
		return mkdirError
	}

	archive, createError := os.Create(zn)
	if createError != nil {
		return createError
	}
	defer func(archive *os.File) { err = archive.Close() }(archive)

	zipWriter := zip.NewWriter(archive)
	defer func(zipWriter *zip.Writer) { err = zipWriter.Close() }(zipWriter)

	entry, entryError := zipWriter.Create(filepath.Base(n))
	if entryError != nil {
		return entryError
	}

	file, openError := os.Open(n)
	if openError != nil {
		return openError
	}

	_, copyError := io.Copy(entry, file)
	if copyError != nil {
		return copyError
	}

	return nil
}

// ZipDirectory zips a directory to the disk.
func ZipDirectory(dn string, zn string) (err error) {
	mkdirError := os.MkdirAll(filepath.Dir(zn), os.ModePerm)
	if mkdirError != nil {
		return mkdirError
	}

	archive, archiveError := os.Create(zn)
	if archiveError != nil {
		return archiveError
	}
	defer func(archive *os.File) { err = archive.Close() }(archive)

	writer := zip.NewWriter(archive)
	defer func(writer *zip.Writer) { err = writer.Close() }(writer)

	return filepath.Walk(dn, func(path string, info fs.FileInfo, err error) error {
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

		entry, entryError := writer.Create(strings.TrimPrefix(path, dn+"/"))
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
func DownloadFile(url string, n string) error {
	response, getError := http.Get(url)
	if getError != nil {
		return getError
	}

	data, readError := io.ReadAll(response.Body)
	if readError != nil {
		return readError
	}

	parentName := filepath.Dir(n)
	if !IsDirectory(parentName) {
		mkdirError := os.MkdirAll(parentName, os.ModePerm)
		if mkdirError != nil {
			return mkdirError
		}
	}

	writeErr := os.WriteFile(n, data, os.ModePerm)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

func FileReader(n string) (*bytes.Reader, os.FileInfo, error) {
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

// ReadFileInChunks reads a file in chunks.
func ReadFileInChunks(n string, c int, cb func(data []byte) error) (err error) {
	file, openError := os.Open(n)
	if openError != nil {
		return openError
	}
	defer func(file *os.File) { err = file.Close() }(file)

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
			callError := cb(buffer[:count-1])
			if callError != nil {
				return callError
			}
		}
		callError := cb(buffer)
		if callError != nil {
			return callError
		}
	}
}
