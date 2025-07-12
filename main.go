package main

import (
	"embed"
	"fmt"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	flag "github.com/spf13/pflag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var FlagHelp = flag.BoolP("help", "h", false, "shows this help document")
var FlagVersion = flag.BoolP("version", "v", false, "shows the binary version and the project version")
var FlagCreateProject = flag.StringP("create-project", "c", "", fmt.Sprintf("creates a frizzante project"))
var FlagCreateJsLibrary = flag.StringP("create-js-library", "l", "", fmt.Sprintf("creates the frizzante library"))

//go:embed version
//go:embed app/frizzante
var efs embed.FS

func main() {
	flag.Parse()

	var version string

	versionData, versionError := efs.ReadFile("version")
	if versionError != nil {
		log.Fatal(versionError)
	}

	version = string(versionData)

	if *FlagVersion {
		println(version)
		os.Exit(0)
	}

	if *FlagHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *FlagCreateProject != "" {
		downloadError := files.DownloadFile("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip", *FlagCreateProject+".zip")
		if downloadError != nil {
			log.Fatal(downloadError)
		}

		unzipError := files.UnzipFile(*FlagCreateProject+".zip", *FlagCreateProject+".tmp")
		if unzipError != nil {
			log.Fatal(unzipError)
		}

		removeError := os.Remove(*FlagCreateProject + ".zip")
		if removeError != nil {
			log.Fatal(removeError)
		}

		renameError := os.Rename(filepath.Join(*FlagCreateProject+".tmp", "frizzante-starter-main"), *FlagCreateProject)
		if renameError != nil {
			log.Fatal(renameError)
		}

		removeAllError := os.RemoveAll(filepath.Join(*FlagCreateProject + ".tmp"))
		if removeAllError != nil {
			log.Fatal(removeAllError)
		}

		os.Exit(0)
	}

	if *FlagCreateJsLibrary != "" {
		efsFileNames, readDirError := embeds.ReadDir(efs, "app/frizzante")
		if readDirError != nil {
			log.Fatal(readDirError)
			return
		}

		if files.IsDirectory(*FlagCreateJsLibrary) {
			removeAllError := os.RemoveAll(*FlagCreateJsLibrary)
			if removeAllError != nil {
				log.Fatal(removeAllError)
				return
			}
		}

		for _, efsFileName := range efsFileNames {
			fileName := fmt.Sprintf("%s/%s", *FlagCreateJsLibrary, strings.TrimPrefix(efsFileName, "app/frizzante"))
			directoryName := filepath.Dir(fileName)

			if !files.IsDirectory(directoryName) {
				mkdirError := os.MkdirAll(directoryName, os.ModePerm)
				if mkdirError != nil {
					log.Fatal(mkdirError)
					return
				}
			}

			file, openError := os.Create(fileName)
			if openError != nil {
				log.Fatal(openError)
				return
			}

			esfFile, esfOpenError := efs.Open(efsFileName)
			if esfOpenError != nil {
				log.Fatal(esfOpenError)
				return
			}

			_, copyError := io.Copy(file, esfFile)
			if copyError != nil {
				log.Fatal(copyError)
				return
			}
		}

		os.Exit(0)
	}

	flag.Usage()
}
