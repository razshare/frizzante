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

type NameMetadata struct {
	Name                         string
	RelativeFileNameTemplate     string
	FullDirectoryNameCamel       string
	FullFileNameCamel            string
	FullFileNameTitle            string
	RelativeFileNameCamel        string
	BaseFileNameNoExtensionTitle string
	BaseDirectoryName            string
	BaseFileNameTitle            string
}

func findNameMetadata(root string, template string, name string, message string) *NameMetadata {
	if !fileExists(root) {
		writeError := os.MkdirAll(root, os.ModePerm)
		if writeError != nil {
			log.Fatal(writeError)
		}
	}

	trimmedName := strings.Trim(name, "\r\n\t ")

	if "" == name {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(message)
		name, _ = reader.ReadString('\n')
		trimmedName = strings.Trim(name, "\r\n\t ")
		if "" == name {
			return findNameMetadata(root, template, name, message)
		}
	}

	fullRoot, fullFileNameRootError := filepath.Abs(root)
	if nil != fullFileNameRootError {
		log.Fatal(fullFileNameRootError)
	}

	metadata := &NameMetadata{Name: trimmedName}

	extensionName := filepath.Ext(template)

	baseFileName := strings.ReplaceAll(
		strings.ReplaceAll(
			trimmedName,
			".",
			string(filepath.Separator),
		),
		"-",
		string(filepath.Separator),
	)

	relativeFileNameTitle := ""
	for _, pageNameChunked := range strings.Split(baseFileName, "_") {
		for _, section := range strings.Split(pageNameChunked, string(filepath.Separator)) {
			relativeFileNameTitle = relativeFileNameTitle + string(filepath.Separator) + strings.ToTitle(section[0:1]) + section[1:]
		}
	}
	relativeFileNameTitle = relativeFileNameTitle + extensionName

	relativeFileNameCamel := ""
	for _, pageNameChunked := range strings.Split(baseFileName, "_") {
		for _, section := range strings.Split(pageNameChunked, string(filepath.Separator)) {
			relativeFileNameCamel = relativeFileNameCamel + string(filepath.Separator) + strings.ToLower(section[0:1]) + section[1:]
		}
	}
	relativeFileNameCamel = relativeFileNameCamel[1:] + extensionName

	fullFileNameCamel := filepath.Join(fullRoot, relativeFileNameCamel)
	fullFileNameTitle := filepath.Join(fullRoot, relativeFileNameTitle)

	metadata.FullDirectoryNameCamel = filepath.Dir(fullFileNameCamel)
	metadata.FullFileNameCamel = fullFileNameCamel
	metadata.FullFileNameTitle = fullFileNameTitle
	metadata.RelativeFileNameCamel = relativeFileNameCamel
	metadata.RelativeFileNameTemplate = template
	metadata.BaseDirectoryName = filepath.Base(metadata.FullDirectoryNameCamel)
	metadata.BaseFileNameTitle = filepath.Base(metadata.FullFileNameTitle)
	metadata.BaseFileNameNoExtensionTitle = strings.TrimSuffix(metadata.BaseFileNameTitle, extensionName)
	return metadata
}

func findNameMetadataForApi(name string) *NameMetadata {
	return findNameMetadata(
		filepath.Join("lib", "api"),
		filepath.Join("templates", "api", "example.go"),
		name,
		"Name the api: ",
	)
}

func findNameMetadataForPage(name string) *NameMetadata {
	return findNameMetadata(
		filepath.Join("lib", "pages"),
		filepath.Join("templates", "pages", "example.go"),
		name,
		"Name the page: ",
	)
}

func findNameMetadataForView(name string) *NameMetadata {
	return findNameMetadata(
		filepath.Join("lib", "components", "views"),
		filepath.Join("templates", "views", "example.svelte"),
		name,
		"Name the view: ",
	)
}

func createApi(apiName string) {
	metadata := findNameMetadataForApi(apiName)

	if !fileExists(metadata.FullDirectoryNameCamel) {
		mkdirError := os.MkdirAll(metadata.FullDirectoryNameCamel, os.ModePerm)
		if mkdirError != nil {
			log.Fatal(mkdirError)
		}
	}

	if fileExists(metadata.FullFileNameCamel) {
		fmt.Printf("Api `%s` already exists.\n", metadata.BaseFileNameNoExtensionTitle)
		createApi("")
		return
	}

	readBytes, readError := templates.ReadFile(strings.ReplaceAll(metadata.RelativeFileNameTemplate, string(filepath.Separator), "/"))
	if nil != readError {
		log.Fatal(readError)
	}

	// Package.
	oldName := []byte("package api")
	newName := []byte("package " + metadata.BaseDirectoryName)
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// Pattern.
	oldName = []byte("\"GET /\"")
	newName = []byte("\"GET /api/" + strings.ReplaceAll(strings.TrimSuffix(metadata.RelativeFileNameCamel, ".go"), string(filepath.Separator), "/") + "\"")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	// Api.
	oldName = []byte("func api(")
	newName = []byte("func " + metadata.BaseFileNameNoExtensionTitle + "(")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	writeError := os.WriteFile(metadata.FullFileNameCamel, readBytes, os.ModePerm)
	if writeError != nil {
		log.Fatal(writeError)
	}
}

func createPage(pageName string) {
	metadata := findNameMetadataForPage(pageName)

	if !fileExists(metadata.FullDirectoryNameCamel) {
		mkdirError := os.MkdirAll(metadata.FullDirectoryNameCamel, os.ModePerm)
		if mkdirError != nil {
			log.Fatal(mkdirError)
		}
	}

	if fileExists(metadata.FullFileNameCamel) {
		fmt.Printf("Page `%s` already exists.\n", metadata.BaseFileNameNoExtensionTitle)
		createPage("")
		return
	}

	readBytes, readError := templates.ReadFile(strings.ReplaceAll(metadata.RelativeFileNameTemplate, string(filepath.Separator), "/"))
	if nil != readError {
		log.Fatal(readError)
	}

	// Package.
	oldName := []byte("package pages")
	newName := []byte("package " + metadata.BaseDirectoryName)
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// PageBuilder.
	oldName = []byte("func page(")
	newName = []byte("func " + metadata.BaseFileNameNoExtensionTitle + "(")
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// View.
	oldName = []byte("\"ViewName\"")
	newName = []byte("\"" + metadata.BaseFileNameNoExtensionTitle + "\"")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	// Path.
	oldName = []byte("\"/path\"")
	newName = []byte("\"/" + strings.ReplaceAll(strings.TrimSuffix(metadata.RelativeFileNameCamel, ".go"), string(filepath.Separator), "/") + "\"")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	writeError := os.WriteFile(metadata.FullFileNameCamel, readBytes, os.ModePerm)
	if writeError != nil {
		log.Fatal(writeError)
	}

	metadata.FullFileNameCamel = strings.TrimSuffix(metadata.FullFileNameCamel, ".go") + ".svelte"
	metadata.RelativeFileNameTemplate = strings.TrimSuffix(metadata.RelativeFileNameTemplate, ".go") + ".svelte"
	metadata.RelativeFileNameCamel = strings.TrimSuffix(metadata.RelativeFileNameCamel, ".go") + ".svelte"
	createViewComponent(metadata.Name)
}

func createViewComponent(pageName string) {
	metadata := findNameMetadataForView(pageName)

	fileName := filepath.Join(metadata.FullDirectoryNameCamel, metadata.BaseFileNameTitle)

	if fileExists(fileName) {
		fmt.Printf("component `%s` already exists.\n", fileName)
		createViewComponent("")
		return
	}

	readBytes, readError := templates.ReadFile(strings.ReplaceAll(metadata.RelativeFileNameTemplate, string(filepath.Separator), "/"))
	if nil != readError {
		log.Fatal(readError)
	}

	// Content.
	oldName := []byte("Hello, this is view!")
	newName := []byte("Hello, this is " + metadata.BaseFileNameTitle + "!")
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	writeError := os.WriteFile(fileName, readBytes, os.ModePerm)
	if writeError != nil {
		log.Fatal(writeError)
	}
}

// Make makes things.
func Make() {
	api := flag.Bool("api", false, "")
	page := flag.Bool("page", false, "")
	name := flag.String("name", "", "")
	flag.Parse()

	if *api {
		createApi(*name)
	}

	if *page {
		createPage(*name)
	}
}
