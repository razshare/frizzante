package cli

import (
	"archive/zip"
	"github.com/pterm/pterm"
	"github.com/razshare/frizzante/fs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Project() {
	if !*FlagGenerate || !*FlagProject {
		return
	}

	var showError error

	if "" == *FlagOut {
		*FlagOut, showError = pterm.
			DefaultInteractiveTextInput.
			Show("The name of the project is")
		if showError != nil {
			log.Fatal(showError)
		}
	}

	response, getError := http.Get("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip")
	if getError != nil {
		log.Fatal(getError)
	}

	zipData, zipReadError := io.ReadAll(response.Body)
	if zipReadError != nil {
		log.Fatal(zipReadError)
	}

	zipFileName := *FlagOut + ".zip"

	zipWriteError := os.WriteFile(zipFileName, zipData, os.ModePerm)
	if zipWriteError != nil {
		log.Fatal(zipWriteError)
	}

	zipReader, zipOpenError := zip.OpenReader(zipFileName)
	if zipOpenError != nil {
		log.Fatal(zipOpenError)
	}
	defer func(zipReader *zip.ReadCloser) {
		err := zipReader.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(zipReader)

	for _, file := range zipReader.File {
		fileName := filepath.Join(*FlagOut, strings.TrimPrefix(file.Name, "frizzante-starter-main"))
		fileIsDirectory := file.FileInfo().IsDir()
		fileIsDirectoryOnDisk := fs.IsDirectory(fileName)

		if fileIsDirectory && !fileIsDirectoryOnDisk {
			mkdirError := os.MkdirAll(fileName, os.ModePerm)
			if mkdirError != nil {
				log.Fatal(mkdirError)
			}
			continue
		}

		directoryName := filepath.Dir(fileName)

		if "." == directoryName {
			continue
		}

		parentExists := fs.IsDirectory(directoryName)

		if !parentExists {
			mkdirError := os.MkdirAll(directoryName, os.ModePerm)
			if mkdirError != nil {
				log.Fatal(mkdirError)
			}
		}

		dstFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := file.Open()
		if err != nil {
			panic(err)
		}

		if _, copyError := io.Copy(dstFile, fileInArchive); copyError != nil {
			panic(copyError)
		}

		_ = dstFile.Close()
		_ = fileInArchive.Close()
	}

	_ = os.Remove(zipFileName)

	os.Exit(0)
}
