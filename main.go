package main

import (
	"embed"
	"fmt"
	"github.com/razshare/frizzante/files"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

var FlagHelp = flag.BoolP("help", "h", false, "shows the help document")
var FlagVersion = flag.BoolP("version", "v", false, "shows the binary version and the project version")
var FlagCreateProject = flag.StringP("create-project", "c", "", fmt.Sprintf("creates a frizzante project"))

//go:embed version
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

		unzipError := files.UnzipFile(*FlagCreateProject+".zip", *FlagCreateProject)
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
