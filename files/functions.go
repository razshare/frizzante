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
func IsFile(fname string) bool {
	stat, statError := os.Stat(fname)
	if statError != nil {
		return false
	}
	return !stat.IsDir()
}

// IsDirectory checks if file exists and is a directory.
func IsDirectory(name string) bool {
	stat, statError := os.Stat(name)
	if statError != nil {
		return false
	}
	return stat.IsDir()
}

// DeleteFile deletes a file from the disk.
func DeleteFile(name string) bool {
	removeError := os.Remove(name)
	return nil == removeError || !errors.Is(removeError, os.ErrNotExist)
}

// UnzipFile unzips a file to the disk.
func UnzipFile(zfname string, dname string) (err error) {
	zipReader, zipOpenError := zip.OpenReader(zfname)
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
		fileNameLocal := filepath.Join(dname, file.Name)
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
func ZipFile(fname string, zfname string) (err error) {
	mkdirError := os.MkdirAll(filepath.Dir(zfname), os.ModePerm)
	if mkdirError != nil {
		return mkdirError
	}

	archive, createError := os.Create(zfname)
	if createError != nil {
		return createError
	}
	defer func(archive *os.File) { err = archive.Close() }(archive)

	writer := zip.NewWriter(archive)
	defer func(writer *zip.Writer) { err = writer.Close() }(writer)

	entry, entryError := writer.Create(filepath.Base(fname))
	if entryError != nil {
		return entryError
	}

	file, fileError := os.Open(fname)
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
func ZipDirectory(dname string, zfname string) (err error) {
	mkdirError := os.MkdirAll(filepath.Dir(zfname), os.ModePerm)
	if mkdirError != nil {
		return mkdirError
	}

	archive, archiveError := os.Create(zfname)
	if archiveError != nil {
		return archiveError
	}
	defer func(archive *os.File) { err = archive.Close() }(archive)

	writer := zip.NewWriter(archive)
	defer func(writer *zip.Writer) { err = writer.Close() }(writer)

	return filepath.Walk(dname, func(path string, info fs.FileInfo, err error) error {
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

		entry, entryError := writer.Create(strings.TrimPrefix(path, dname+"/"))
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
func DownloadFile(url string, fname string) error {
	response, getError := http.Get(url)
	if getError != nil {
		return getError
	}

	data, readError := io.ReadAll(response.Body)
	if readError != nil {
		return readError
	}

	parentName := filepath.Dir(fname)
	if !IsDirectory(parentName) {
		mkdirError := os.MkdirAll(parentName, os.ModePerm)
		if mkdirError != nil {
			return mkdirError
		}
	}

	writeErr := os.WriteFile(fname, data, os.ModePerm)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

func FileReader(fname string) (*bytes.Reader, os.FileInfo, error) {
	file, openError := os.Open(fname)
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
func ReadFileInChunks(fname string, chunk int, fun func(data []byte) error) (err error) {
	file, openError := os.Open(fname)
	if openError != nil {
		return openError
	}
	defer func(file *os.File) { err = file.Close() }(file)

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
			callError := fun(buffer[:count-1])
			if callError != nil {
				return callError
			}
		}
		callError := fun(buffer)
		if callError != nil {
			return callError
		}
	}
}
