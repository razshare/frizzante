package frizzante

import (
	"bufio"
	"bytes"
	"embed"
	"flag"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*/**
var templates embed.FS

type CliMetadataPayload struct {
	Root     string
	Template string
	Name     string
	Base     string
	Message  string
}

func template(from string, to string, replacements map[string]string) {
	if fileExists(to) {
		fmt.Printf("file `%s` already exists\n", to)
		return
	}

	readBytes, readError := templates.ReadFile(from)
	if nil != readError {
		log.Fatal(readError)
	}

	for key, replacement := range replacements {
		readBytes = bytes.ReplaceAll(readBytes, []byte(key), []byte(replacement))
	}

	directory := filepath.Dir(to)

	if !fileExists(directory) {
		err := os.MkdirAll(directory, os.ModePerm)
		if nil != err {
			log.Fatal(err)
		}
	}
	err := os.WriteFile(to, readBytes, os.ModePerm)
	if nil != err {
		log.Fatal(err)
	}
}

// Cli makes things.
func Cli() {
	titleCase := cases.Title(language.English)
	route := flag.Bool("route", false, "")
	name := flag.String("name", "", "")
	flag.Parse()

	*name = strings.Trim(*name, "\r\n\t ")

	if *route {
		for "" == *name {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("What's the name of the route? ")
			*name, _ = reader.ReadString('\n')
		}

		fixedName := strings.Trim(*name, "\r\n\t ")
		packageName := fixedName
		fileNameGo := fixedName + ".go"
		fileNameSvelte := titleCase.String(fixedName) + ".svelte"
		functionName := "Get" + titleCase.String(fixedName)
		viewName := titleCase.String(fixedName)

		template(
			filepath.Join("templates", "routes", "example.go"),
			filepath.Join("lib", "routes", fileNameGo),
			map[string]string{
				"package pages": fmt.Sprintf("package %s", packageName),
				"functionName":  functionName,
				"viewName":      viewName,
			},
		)

		template(
			filepath.Join("templates", "views", "example.svelte"),
			filepath.Join("lib", "components", "views", fileNameSvelte),
			map[string]string{},
		)
	}
}
