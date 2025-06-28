package main

import (
	"embed"
	"fmt"
	ffs "github.com/razshare/frizzante/fs"
	flag "github.com/spf13/pflag"
	"io"
	"io/fs"
	"log"
	"os"
)

const binaryVersion = "v1.8.2"

var FlagHelp = flag.BoolP("help", "h", false, "shows the help document")
var FlagVersion = flag.BoolP("version", "v", false, "shows the binary version and the project version")
var FlagCreateProject = flag.StringP("create-project", "c", "", fmt.Sprintf("creates a frizzante project"))
var FlagRestoreUtilities = flag.StringP("restore-utilities", "r", "", fmt.Sprintf("restores frizzante utilities"))

//go:embed app/lib/utilities
//go:embed project.zip
var embedded embed.FS

func main() {
	flag.Parse()

	if *FlagRestoreUtilities != "" {
		var restore func(from string, to string)
		restore = func(from string, to string) {
			if !ffs.IsDirectory(to) {
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
		println(binaryVersion)
		os.Exit(0)
	}

	if *FlagHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *FlagCreateProject != "" {
		srcFile, srcError := embedded.Open("project.zip")
		if srcError != nil {
			log.Fatal(srcError)
		}
		defer func(srcFile fs.File) {
			err := srcFile.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(srcFile)

		destFile, destError := os.Create(*FlagCreateProject + ".zip")
		if destError != nil {
			log.Fatal(destError)
		}
		defer func(destFile *os.File) {
			err := destFile.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(destFile)

		_, copyError := io.Copy(destFile, srcFile)
		if copyError != nil {
			return
		}

		unzipError := ffs.UnzipFile(*FlagCreateProject+".zip", *FlagCreateProject)
		if unzipError != nil {
			log.Fatal(unzipError)
		}

		removeError := os.Remove(*FlagCreateProject + ".zip")
		if removeError != nil {
			log.Fatal(removeError)
		}
		os.Exit(0)
	}

	flag.Usage()
}
