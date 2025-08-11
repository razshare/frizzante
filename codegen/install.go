package codegen

import (
	"fmt"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"path/filepath"
	"strings"
)

func Install(name string, url string, destination string) {
	if files.IsDirectory(destination) {
		if !confirm.Sendf(true, "It looks like `%s` is already installed in `%s`, would you like to overwrite it?", name, destination) {
			messages.Infof("skipping `%s`", name)
			return
		}

		removeError := os.RemoveAll(destination)
		if removeError != nil {
			messages.Fatal(removeError)
		}
	}

	spin := spinner.New(fmt.Sprintf("installing `%s` from `%s`...", name, url))
	spinnerError := spinner.Start(spin)
	if spinnerError != nil {
		messages.Fatal(spinnerError)
	}
	defer func() { spinner.Stop(spin) }()

	if !strings.HasSuffix(url, ".zip") {
		nameFixed := name
		if strings.HasSuffix(url, ".exe") {
			nameFixed += ".exe"
		}

		downloadError := files.DownloadFile(url, filepath.Join(destination, nameFixed))
		if downloadError != nil {
			messages.Fatal(downloadError)
		}

		messages.Successf("%s installed in `%s`", name, destination)
		return
	}

	zipFileName := destination + ".zip"
	downloadError := files.DownloadFile(url, zipFileName)
	if downloadError != nil {
		messages.Fatal(downloadError)
	}
	defer func() {
		removeError := os.Remove(zipFileName)
		if removeError != nil {
			messages.Fatal(removeError)
		}
	}()

	unzipError := files.UnzipFile(zipFileName, destination)
	if unzipError != nil {
		messages.Fatal(unzipError)
	}

	messages.Successf("%s installed in `%s`", name, destination)
}
