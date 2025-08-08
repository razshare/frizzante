package cli

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/razshare/frizzante/files"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func OnInstall() {
	OnTouch()

	tidy := exec.Command(Go("."), "mod", "tidy")
	tidy.Env = append(os.Environ())
	tidy.Stderr = os.Stderr
	tidy.Stdout = os.Stdout
	tidy.Stdin = os.Stdin
	tidyError := tidy.Run()
	if tidyError != nil {
		Fatal(tidyError)
	}

	install := exec.Command(Bun("app"), "install")
	install.Dir = "app"
	install.Env = append(os.Environ())
	install.Stderr = os.Stderr
	install.Stdout = os.Stdout
	install.Stdin = os.Stdin
	installError := install.Run()
	if installError != nil {
		Fatal(installError)
	}

	Success("project dependencies installed")
}

func Install(name string, url string, destination string) {
	if files.IsDirectory(destination) {
		if !Confirmf("It looks like `%s` is already installed in `%s`, would you like to overwrite it?", name, destination) {
			Infof("skipping `%s`", name)
			return
		}

		removeError := os.RemoveAll(destination)
		if removeError != nil {
			Fatal(removeError)
		}
	}

	spinner, spinnerError := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start(fmt.Sprintf("installing `%s` from `%s`...", name, url))
	if spinnerError != nil {
		Fatal(spinnerError)
	}
	defer func() {
		stopError := spinner.Stop()
		if stopError != nil {
			Fatal(stopError)
		}
	}()

	if !strings.HasSuffix(url, ".zip") {
		nameFixed := name
		if strings.HasSuffix(url, ".exe") {
			nameFixed += Extension()
		}

		downloadError := files.DownloadFile(url, filepath.Join(destination, nameFixed))
		if downloadError != nil {
			Fatal(downloadError)
		}

		Successf("%s installed in `%s`", name, destination)
		return
	}

	zipFileName := destination + ".zip"
	downloadError := files.DownloadFile(url, zipFileName)
	if downloadError != nil {
		Fatal(downloadError)
	}
	defer func() {
		removeError := os.Remove(zipFileName)
		if removeError != nil {
			Fatal(removeError)
		}
	}()

	unzipError := files.UnzipFile(zipFileName, destination)
	if unzipError != nil {
		Fatal(unzipError)
	}

	Successf("%s installed in `%s`", name, destination)
}
