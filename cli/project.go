package cli

import (
	"archive/zip"
	"github.com/pterm/pterm"
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

	for _, file := range zipReader.File {
		fileName := filepath.Join(*FlagOut, strings.TrimPrefix(file.Name, "frizzante-starter-main"))
		fileIsDirectory := file.FileInfo().IsDir()
		fileIsDirectoryOnDisk := isDirectory(fileName)

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

		parentExists := isDirectory(directoryName)

		if !parentExists {
			mkdirError := os.MkdirAll(directoryName, os.ModePerm)
			if mkdirError != nil {
				log.Fatal(mkdirError)
			}
		}

		fileReader, openRawError := file.OpenRaw()
		if openRawError != nil {
			log.Fatal(openRawError)
		}

		data, readError := io.ReadAll(fileReader)
		if readError != nil {
			log.Fatal(readError)
		}

		writeError := os.WriteFile(fileName, data, os.ModePerm)
		if writeError != nil {
			log.Fatal(writeError)
		}
	}

	_ = os.Remove(zipFileName)

	os.Exit(0)
}
