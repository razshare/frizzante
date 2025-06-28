package main

import (
	"embed"
	"fmt"
	"github.com/razshare/frizzante/fs"
	flag "github.com/spf13/pflag"
	"log"
	"os"
	"path/filepath"
)

const binaryVersion = "1.8.3"
const projectVersion = "1.2.7"

var FlagHelp = flag.BoolP("help", "h", false, "shows the help document")
var FlagVersion = flag.BoolP("version", "v", false, "shows the binary version and the project version")
var FlagCreateProject = flag.StringP("create-project", "c", "", fmt.Sprintf("creates a frizzante project"))
var FlagRestoreUtilities = flag.StringP("restore-utilities", "r", "", fmt.Sprintf("restores frizzante utilities"))

//go:embed app/lib/utilities
var embedded embed.FS

func main() {
	flag.Parse()

	if *FlagRestoreUtilities != "" {
		var restore func(from string, to string)
		restore = func(from string, to string) {
			if !fs.IsDirectory(to) {
				mkdirAllError := os.MkdirAll(to, os.ModePerm)
				if mkdirAllError != nil {
					log.Fatal(mkdirAllError)
				}
			}

			entries, readDirError := embedded.ReadDir(from)
			if readDirError != nil {
				log.Fatal(readDirError)
			}

			for _, entry := range entries {
				if entry.IsDir() {
					fromDirectoryName := from + "/" + entry.Name()
					toDirectoryName := to + "/" + entry.Name()
					restore(fromDirectoryName, toDirectoryName)
					continue
				}

				fromFileName := from + "/" + entry.Name()
				toFileName := to + "/" + entry.Name()

				data, readError := embedded.ReadFile(fromFileName)
				if readError != nil {
					log.Fatal(readError)
				}

				writeError := os.WriteFile(toFileName, data, os.ModePerm)
				if writeError != nil {
					log.Fatal(readError)
				}
			}
		}

		restore("app/lib/utilities", *FlagRestoreUtilities)
		os.Exit(0)
	}

	if *FlagVersion {
		fmt.Printf("v%s using project v%s", binaryVersion, projectVersion)
		os.Exit(0)
	}

	if *FlagHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *FlagCreateProject != "" {
		downloadError := fs.DownloadFile(fmt.Sprintf("https://github.com/razshare/frizzante-starter/archive/refs/tags/v%s.zip", projectVersion), *FlagCreateProject+".zip")
		if downloadError != nil {
			log.Fatal(downloadError)
		}

		unzipError := fs.UnzipFile(*FlagCreateProject+".zip", "."+*FlagCreateProject+".tmp")
		if unzipError != nil {
			log.Fatal(unzipError)
		}

		removeError := os.Remove(*FlagCreateProject + ".zip")
		if removeError != nil {
			log.Fatal(removeError)
		}

		renameError := os.Rename(filepath.Join("."+*FlagCreateProject+".tmp", fmt.Sprintf("frizzante-starter-%s", projectVersion)), *FlagCreateProject)
		if renameError != nil {
			log.Fatal(renameError)
		}

		removeError = os.RemoveAll("." + *FlagCreateProject + ".tmp")
		if removeError != nil {
			log.Fatal(removeError)
		}
		os.Exit(0)
	}

	flag.Usage()
}
