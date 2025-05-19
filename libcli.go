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
	"unicode"
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
	RelativeFileNameTitle        string
	BaseFileNameNoExtensionTitle string
	BaseDirectoryName            string
	BaseFileNameTitle            string
	BaseFileNameKebab            string
}

func toKebab(value string) string {
	result := ""
	for _, item := range value {
		if unicode.IsUpper(item) {
			result += "-" + string(unicode.ToLower(item))
			continue
		}
		if '_' == item || '.' == item {
			result += "-"
		}
		result += string(item)
	}
	return strings.Trim(result, "-")
}

func cleanUpName(name string) string {
	name = strings.Trim(name, "\r\n\t ")
	name = strings.TrimPrefix(name, "Api")
	name = strings.TrimSuffix(name, "Api")
	name = strings.TrimPrefix(name, "api")
	name = strings.TrimSuffix(name, "api")
	name = strings.TrimPrefix(name, "Controller")
	name = strings.TrimSuffix(name, "Controller")
	name = strings.TrimPrefix(name, "controller")
	name = strings.TrimSuffix(name, "controller")
	name = strings.TrimPrefix(name, "View")
	name = strings.TrimSuffix(name, "View")
	name = strings.TrimPrefix(name, "view")
	name = strings.TrimSuffix(name, "view")
	return name
}

func findNameMetadata(root string, template string, name string, message string) *NameMetadata {
	if !fileExists(root) {
		writeError := os.MkdirAll(root, os.ModePerm)
		if writeError != nil {
			log.Fatal(writeError)
		}
	}

	name = cleanUpName(name)

	if "" == name {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(message)
		name, _ = reader.ReadString('\n')
		if "" == name {
			return findNameMetadata(root, template, name, message)
		}
	}

	name = cleanUpName(name)

	fullRoot, fullFileNameRootError := filepath.Abs(root)
	if nil != fullFileNameRootError {
		log.Fatal(fullFileNameRootError)
	}

	metadata := &NameMetadata{Name: name}

	extensionName := filepath.Ext(template)

	baseFileName := strings.ReplaceAll(
		strings.ReplaceAll(
			name,
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
	relativeFileNameTitle = relativeFileNameTitle[1:] + extensionName

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
	metadata.RelativeFileNameTitle = relativeFileNameTitle
	metadata.RelativeFileNameTemplate = template
	metadata.BaseDirectoryName = filepath.Base(metadata.FullDirectoryNameCamel)
	metadata.BaseFileNameTitle = filepath.Base(metadata.FullFileNameTitle)
	metadata.BaseFileNameNoExtensionTitle = strings.TrimSuffix(metadata.BaseFileNameTitle, extensionName)
	metadata.BaseFileNameKebab = toKebab(metadata.BaseFileNameNoExtensionTitle)
	return metadata
}

func findNameMetadataForApi(name string) *NameMetadata {
	return findNameMetadata(
		filepath.Join("lib", "controllers", "api"),
		filepath.Join("templates", "api", "example.go"),
		name,
		"What's the name of this api? ",
	)
}

func findNameMetadataForPage(name string) *NameMetadata {
	return findNameMetadata(
		filepath.Join("lib", "controllers", "pages"),
		filepath.Join("templates", "pages", "example.go"),
		name,
		"What's the name of this page? ",
	)
}

func findNameMetadataForView(name string) *NameMetadata {
	return findNameMetadata(
		filepath.Join("lib", "components", "views"),
		filepath.Join("templates", "views", "example.svelte"),
		name,
		"What is the name of this view? ",
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

	if fileExists(metadata.FullFileNameTitle) {
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
	oldName = []byte("/path")
	newName = []byte("/api/" + strings.ReplaceAll(strings.TrimSuffix(metadata.BaseFileNameKebab, ".go"), string(filepath.Separator), "/"))
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// Api.
	oldName = []byte("apiController")
	newName = []byte(metadata.BaseFileNameNoExtensionTitle + "Controller")
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	writeError := os.WriteFile(strings.TrimSuffix(metadata.FullFileNameTitle, ".go")+"Controller.go", readBytes, os.ModePerm)
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

	fileName := strings.TrimSuffix(metadata.FullFileNameTitle, ".go") + "Controller.go"

	if !fileExists(metadata.FullFileNameTitle) {
		readBytes, readError := templates.ReadFile(strings.ReplaceAll(metadata.RelativeFileNameTemplate, string(filepath.Separator), "/"))
		if nil != readError {
			log.Fatal(readError)
		}

		// Package.
		oldName := []byte("package pages")
		newName := []byte("package " + metadata.BaseDirectoryName)
		readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

		// PageController.
		oldName = []byte("pageController")
		newName = []byte(metadata.BaseFileNameNoExtensionTitle + "Controller")
		readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

		// PageData.
		oldName = []byte("pageData")
		newName = []byte(metadata.BaseFileNameNoExtensionTitle + "Data")
		readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

		// Path.
		oldName = []byte("/path")
		newName = []byte("/" + strings.ReplaceAll(strings.TrimSuffix(metadata.BaseFileNameKebab, ".go"), string(filepath.Separator), "/"))
		readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

		writeError := os.WriteFile(fileName, readBytes, os.ModePerm)
		if writeError != nil {
			log.Fatal(writeError)
		}
	} else {
		fmt.Printf("Page `%s` already exists.\n", metadata.BaseFileNameNoExtensionTitle)
	}

	metadata.FullFileNameTitle = strings.TrimSuffix(metadata.FullFileNameTitle, ".go") + ".svelte"
	metadata.RelativeFileNameTemplate = strings.TrimSuffix(metadata.RelativeFileNameTemplate, ".go") + ".svelte"
	metadata.RelativeFileNameCamel = strings.TrimSuffix(metadata.RelativeFileNameCamel, ".go") + ".svelte"
	createViewComponent(metadata.Name)
}

func createViewComponent(pageName string) {
	metadata := findNameMetadataForView(pageName)
	fileName := strings.TrimSuffix(metadata.FullFileNameTitle, ".svelte") + "View.svelte"

	if !fileExists(metadata.FullDirectoryNameCamel) {
		mkdirError := os.MkdirAll(metadata.FullDirectoryNameCamel, os.ModePerm)
		if mkdirError != nil {
			log.Fatal(mkdirError)
		}
	}

	if !fileExists(metadata.FullFileNameTitle) {
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
	} else {
		fmt.Printf("component `%s` already exists.\n", metadata.BaseFileNameNoExtensionTitle)
	}
}

// Cli makes things.
func Cli() {
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
