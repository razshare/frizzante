package cli

import (
	"github.com/pterm/pterm"
	"github.com/razshare/frizzante/fs"
	"io"
	"log"
	"net/http"
	"os"
)

func Project() {
	if !*FlagCreateProject {
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

	unzipError := fs.UnzipFile(zipFileName, *FlagOut)
	if unzipError != nil {
		log.Fatal(unzipError)
	}
}
