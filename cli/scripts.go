package cli

import (
	"embed"
	"flag"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/razshare/frizzante/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

var FlagDevelop = flag.Bool("develop", false, "")
var FlagCreateProject = flag.Bool("project", false, "")
var FlagGenerateUtilities = flag.Bool("utilities", false, "")
var FlagOut = flag.String("out", "", "")
var FlagStarterVersion = flag.String("starter-version", "1.2.5", "")
var FlagAirVersion = flag.String("air-version", "1.62.0", "")
var FlagBunVersion = flag.String("bun-version", "1.2.16", "")

//go:embed .air.toml
//go:embed utilities
var efs embed.FS

func GenerateUtilities() {
	if !*FlagGenerateUtilities {
		return
	}

	var showError error

	if "" == *FlagOut {
		*FlagOut, showError = pterm.
			DefaultInteractiveTextInput.
			Show("Drop utilities in")
		if showError != nil {
			log.Fatal(showError)
		}
	}

	u := NewUtilities(efs)

	if showError = u.CreateOnDisk(filepath.Join(*FlagOut)); showError != nil {
		log.Fatal(showError)
	}

	os.Exit(0)
}

func CreateProject() {
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

	fileName := *FlagOut + ".zip"
	directoryName := *FlagOut
	directoryNameTemp := *FlagOut + ".tmp"

	if fs.FileExists(directoryName) {
		log.Fatal(fmt.Sprintf("%s already exists", directoryName))
	}

	if fs.FileExists(directoryNameTemp) {
		log.Fatal(fmt.Sprintf("%s already exists", directoryNameTemp))
	}

	if fs.FileExists(fileName) {
		log.Fatal(fmt.Sprintf("%s already exists", fileName))
	}

	home, homeError := os.UserHomeDir()
	if homeError != nil {
		log.Fatal(homeError)
	}

	base := fmt.Sprintf("%s/.frizzante", home)

	starterUrl := "https://github.com/razshare/frizzante-starter/archive/refs/tags/v" + *FlagStarterVersion + ".zip"

	if !fs.FileExists(fmt.Sprintf("%s/starter.zip", base)) {
		downloadError := fs.DownloadFile(starterUrl, fmt.Sprintf("%s/starter.zip", base))
		if downloadError != nil {
			log.Fatal(downloadError)
		}
	}

	unzipError := fs.UnzipFile(fmt.Sprintf("%s/starter.zip", base), directoryNameTemp)
	if unzipError != nil {
		log.Fatal(unzipError)
	}

	renameError := os.Rename(filepath.Join(directoryNameTemp, "frizzante-starter-"+*FlagStarterVersion), directoryName)
	if renameError != nil {
		log.Fatal(renameError)
	}

	removeError := os.RemoveAll(directoryNameTemp)
	if removeError != nil {
		log.Fatal(removeError)
	}
}

func Develop() {
	if !*FlagDevelop {
		return
	}

	var platform string
	var bunUrl string
	var airUrl string
	var extension string

	if runtime.GOOS == "linux" {
		platform = "linux"
		bunUrl = "https://github.com/oven-sh/bun/releases/download/bun-v" + *FlagBunVersion + "/bun-linux-x64.zip"
		airUrl = "https://github.com/air-verse/air/releases/download/v" + *FlagAirVersion + "/air_" + *FlagAirVersion + "_linux_amd64"
	} else if runtime.GOOS == "darwin" {
		platform = "darwin"
		bunUrl = "https://github.com/oven-sh/bun/releases/download/bun-v" + *FlagBunVersion + "/bun-darwin-x64.zip"
		airUrl = "https://github.com/air-verse/air/releases/download/v" + *FlagAirVersion + "/air_" + *FlagAirVersion + "_darwin_amd64"
	} else if runtime.GOOS == "windows" {
		platform = "windows"
		extension = ".exe"
		bunUrl = "https://github.com/oven-sh/bun/releases/download/bun-v" + *FlagBunVersion + "6/bun-windows-x64.zip"
		airUrl = "https://github.com/air-verse/air/releases/download/v" + *FlagAirVersion + "/air_" + *FlagAirVersion + "_windows_amd64" + extension
	} else {
		log.Fatalf("unknown platform %s", runtime.GOOS)
	}

	if filepath.Separator == '\\' {
		platform = "windows"
	}

	if !fs.FileExists("main.go") {
		log.Fatal("no main.go file detected")
	}

	if !fs.FileExists(".air.toml") {
		data, readError := efs.ReadFile(".air.toml")
		if readError != nil {
			log.Fatal(readError)
		}

		writeError := os.WriteFile(".air.toml", data, os.ModePerm)
		if writeError != nil {
			log.Fatal(writeError)
		}
	}

	home, homeError := os.UserHomeDir()
	if homeError != nil {
		log.Fatal(homeError)
	}

	base := fmt.Sprintf("%s/.frizzante", home)

	mkdirError := os.MkdirAll(base, os.ModePerm)
	if mkdirError != nil {
		log.Fatal(mkdirError)
	}

	mkdirError = os.MkdirAll("bin", os.ModePerm)
	if mkdirError != nil {
		log.Fatal(mkdirError)
	}

	if !fs.FileExists(fmt.Sprintf("%s/air-%s-%s.zip", base, *FlagAirVersion, platform)) {
		downloadError := fs.DownloadFile(airUrl, fmt.Sprintf("%s/air/air-%s-%s-x64/air%s", base, *FlagAirVersion, platform, extension))
		if downloadError != nil {
			log.Fatal(downloadError)
		}

		zipError := fs.ZipDirectory(fmt.Sprintf("%s/air", base), fmt.Sprintf("%s/air-%s-%s.zip", base, *FlagAirVersion, platform))
		if zipError != nil {
			log.Fatal(zipError)
		}

		removeError := os.RemoveAll(fmt.Sprintf("%s/air", base))
		if removeError != nil {
			log.Fatal(removeError)
		}
	}

	if !fs.FileExists(filepath.Join("bin", "air"+extension)) {
		// Configure air...
		spinner, _ := pterm.DefaultSpinner.
			WithShowTimer(false).
			WithRemoveWhenDone(true).
			Start("Configuring Air...")

		removeError := os.RemoveAll(filepath.Join("bin", "air.tmp"))
		if removeError != nil {
			log.Fatal(removeError)
		}

		unzipError := fs.UnzipFile(fmt.Sprintf("%s/air-%s-%s.zip", base, *FlagAirVersion, platform), filepath.Join("bin", "air.tmp"))
		if unzipError != nil {
			log.Fatal(unzipError)
		}

		renameError := os.Rename(
			filepath.Join("bin", "air.tmp", pterm.Sprintf("air-%s-%s-x64", *FlagAirVersion, platform), "air"+extension),
			filepath.Join("bin", "air"+extension),
		)
		if renameError != nil {
			log.Fatal(renameError)
		}

		removeError = os.RemoveAll(filepath.Join("bin", "air.tmp"))
		if removeError != nil {
			log.Fatal(removeError)
		}

		stopError := spinner.Stop()
		if stopError != nil {
			log.Fatal(stopError)
		}
	}

	if !fs.FileExists(fmt.Sprintf("%s/bun-%s-%s.zip", base, *FlagBunVersion, platform)) {
		downloadError := fs.DownloadFile(bunUrl, fmt.Sprintf("%s/bun-%s-%s.zip", base, *FlagBunVersion, platform))
		if downloadError != nil {
			log.Fatal(downloadError)
		}
	}

	if !fs.FileExists(filepath.Join("bin", "bun"+extension)) {
		// Configure bun...
		spinner, _ := pterm.DefaultSpinner.
			WithShowTimer(false).
			WithRemoveWhenDone(true).
			Start("Configuring Bun...")

		removeError := os.RemoveAll(filepath.Join("bin", "bun.tmp"))
		if removeError != nil {
			log.Fatal(removeError)
		}

		unzipError := fs.UnzipFile(fmt.Sprintf("%s/bun-%s-%s.zip", base, *FlagBunVersion, platform), filepath.Join("bin", "bun.tmp"))
		if unzipError != nil {
			log.Fatal(unzipError)
		}

		renameError := os.Rename(
			filepath.Join("bin", "bun.tmp", pterm.Sprintf("bun-%s-x64", platform), "bun"+extension),
			filepath.Join("bin", "bun"+extension),
		)
		if renameError != nil {
			log.Fatal(renameError)
		}

		removeError = os.RemoveAll(filepath.Join("bin", "bun.tmp"))
		if removeError != nil {
			log.Fatal(removeError)
		}

		stopError := spinner.Stop()
		if stopError != nil {
			log.Fatal(stopError)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	directory, homeError := os.Getwd()
	if homeError != nil {
		log.Println(homeError)
	}

	// Create commands...
	///// bun update
	bunUpdate := exec.Command(filepath.Join("..", "bin", "bun"+extension), "update")
	bunUpdate.Dir = filepath.Join(directory, "app")
	bunUpdate.Stdout = os.Stdout
	bunUpdate.Stderr = os.Stderr
	bunUpdate.Env = append(
		os.Environ(),
		"DEV=1",
	)
	///// tidy
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = filepath.Join(directory, "app")
	tidy.Stdout = os.Stdout
	tidy.Stderr = os.Stderr
	tidy.Env = append(
		os.Environ(),
		"DEV=1",
	)
	///// csr
	csr := exec.Command(
		filepath.Join("..", "bin", "bun"+extension),
		"x",
		"--bun",
		"vite",
		"build",
		"--outDir=dist/client",
		"--logLevel=info",
	)
	csr.Dir = filepath.Join(directory, "app")
	csr.Stdout = os.Stdout
	csr.Stderr = os.Stderr
	csr.Env = append(
		os.Environ(),
		"DEV=1",
	)
	///// ssr
	ssr := exec.Command(
		filepath.Join("..", "bin", "bun"+extension),
		"x",
		"--bun",
		"vite",
		"build",
		"--watch",
		"--outDir=dist",
		"--logLevel=info",
		"--ssr=lib/utilities/frz/scripts/server.ts",
	)
	ssr.Dir = filepath.Join(directory, "app")
	ssr.Stdout = os.Stdout
	ssr.Stderr = os.Stderr
	ssr.Env = append(
		os.Environ(),
		"DEV=1",
	)
	///// air
	air := exec.Command(filepath.Join("bin", "air"+extension))
	air.Dir = directory
	air.Stdout = os.Stdout
	air.Stderr = os.Stderr
	air.Env = append(
		os.Environ(),
		"DEV=1",
		"CGO_ENABLED=1",
	)

	// Run bun update once...
	bunUpdateError := bunUpdate.Run()
	if bunUpdateError != nil {
		log.Fatal(bunUpdateError)
	}

	// Run go mod tidy once...
	tidyError := tidy.Run()
	if tidyError != nil {
		log.Fatal(tidyError)
	}

	// Run vite csr build once...
	csrError := csr.Run()
	if csrError != nil {
		log.Fatal(csrError)
	}

	go func() {
		// Start vite ssr...
		ssrError := ssr.Start()
		if ssrError != nil {
			_ = air.Cancel()
			log.Fatal(ssrError)
		}
	}()

	go func() {
		// Start air...
		airError := air.Start()
		if airError != nil {
			_ = ssr.Cancel()
			log.Fatal(airError)
		}
	}()

	<-sigs

	_ = ssr.Wait()
	_ = air.Wait()

	os.Exit(0)
}
