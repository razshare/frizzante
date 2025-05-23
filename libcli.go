package frizzante

import (
	"bufio"
	"bytes"
	"embed"
	"flag"
	"fmt"
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
	api := flag.Bool("api", false, "")
	page := flag.Bool("page", false, "")
	name := flag.String("name", "", "")
	flag.Parse()

	*name = strings.Trim(*name, "\r\n\t ")

	if *api {
		for "" == *name {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("What's the name of this api? ")
			*name, _ = reader.ReadString('\n')
		}

		*name = strings.Trim(*name, "\r\n\t ")

		template(
			filepath.Join("templates", "api", "example.go"),
			filepath.Join("lib", "api", *name+".go"),
			map[string]string{
				"var Server *f.Server\n": "",
				"//config.Server":        "\"main/lib/config\"",
				"Server.":                "config.Server.",
				"/path":                  fmt.Sprintf("/api/%s", *name),
			},
		)
	}

	if *page {
		for "" == *name {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("What's the name of this page? ")
			*name, _ = reader.ReadString('\n')
		}

		*name = strings.Trim(*name, "\r\n\t ")

		template(
			filepath.Join("templates", "pages", "example.go"),
			filepath.Join("lib", "controllers", *name, "page.go"),
			map[string]string{
				"package pages":          fmt.Sprintf("package %s", *name),
				"var Server *f.Server\n": "",
				"//config.Server":        "\"main/lib/config\"",
				"Server.":                "config.Server.",
			},
		)

		template(
			filepath.Join("templates", "views", "example.svelte"),
			filepath.Join("lib", "controllers", *name, "view.svelte"),
			map[string]string{
				"this is view": fmt.Sprintf("this is %s", *name),
			},
		)
	}
}
