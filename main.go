package main

import (
	"embed"
	"flag"
	"github.com/razshare/frizzante/fs"
	"log"
	"os"
)

var FlagRestoreUtilities = flag.Bool("restore-utilities", false, "")
var FlagOut = flag.String("out", "", "")

//go:embed app/lib/utilities
var utilities embed.FS

func main() {
	flag.Parse()

	if *FlagRestoreUtilities {
		if *FlagOut == "" {
			log.Fatal("flag -out is required")
		}

		var restore func(from string, to string)
		restore = func(from string, to string) {
			if !fs.IsDirectory(to) {
				mkdirAllError := os.MkdirAll(to, os.ModePerm)
				if mkdirAllError != nil {
					log.Fatal(mkdirAllError)
				}
			}

			entries, readDirError := utilities.ReadDir(from)
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

				data, readError := utilities.ReadFile(fromFileName)
				if readError != nil {
					log.Fatal(readError)
				}

				writeError := os.WriteFile(toFileName, data, os.ModePerm)
				if writeError != nil {
					log.Fatal(readError)
				}
			}
		}

		restore("app/lib/utilities", *FlagOut)
	}
}
