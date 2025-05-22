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
	FullFileName                 string
	RelativeFileNameCamel        string
	RelativeFileNameTitle        string
	BaseFileNameNoExtensionTitle string
	RelativeDirectoryNameCamel   string
	RelativeDirectoryNameTitle   string
	BaseDirectoryNameCamel       string
	BaseDirectoryNameTitle       string
	BaseFileNameTitle            string
	BaseFileNameCamel            string
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

func findNameMetadata(root string, template string, name string, base string, message string) *NameMetadata {
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
			return findNameMetadata(root, template, name, base, message)
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
	fullFileName := filepath.Join(fullRoot, baseFileName)

	metadata.FullDirectoryNameCamel = filepath.Dir(fullFileNameCamel)
	metadata.FullFileNameCamel = fullFileNameCamel
	metadata.FullFileNameTitle = fullFileNameTitle
	metadata.FullFileName = fullFileName
	metadata.RelativeFileNameCamel = relativeFileNameCamel
	metadata.RelativeFileNameTitle = relativeFileNameTitle
	metadata.RelativeFileNameTemplate = template
	metadata.RelativeDirectoryNameCamel = filepath.Dir(relativeFileNameCamel)
	metadata.RelativeDirectoryNameTitle = filepath.Dir(relativeFileNameTitle)
	metadata.BaseDirectoryNameCamel = filepath.Base(filepath.Dir(metadata.FullFileNameTitle))
	metadata.BaseDirectoryNameTitle = filepath.Base(filepath.Dir(metadata.FullFileNameCamel))
	metadata.BaseFileNameTitle = filepath.Base(metadata.FullFileNameTitle)
	metadata.BaseFileNameCamel = filepath.Base(metadata.FullFileNameCamel)
	metadata.BaseFileNameNoExtensionTitle = strings.TrimSuffix(metadata.BaseFileNameTitle, extensionName)
	metadata.BaseFileNameKebab = toKebab(metadata.BaseFileNameNoExtensionTitle)
	return metadata
}

func findNameMetadataForApi(name string) *NameMetadata {
	return findNameMetadata(
		filepath.Join("lib", "controllers", "api"),
		filepath.Join("templates", "api", "example.go"),
		name,
		"controller",
		"What's the name of this api? ",
	)
}

func findNameMetadataForPage(name string) *NameMetadata {
	return findNameMetadata(
		filepath.Join("lib", "controllers", "pages"),
		filepath.Join("templates", "pages", "example.go"),
		name,
		"controller",
		"What's the name of this page? ",
	)
}

func findNameMetadataForView(name string) *NameMetadata {
	return findNameMetadata(
		filepath.Join("lib", "controllers", "pages"),
		filepath.Join("templates", "views", "example.svelte"),
		name,
		"view",
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

	fileName := metadata.FullFileName

	if fileExists(fileName) {
		fmt.Printf("file `%s` already exists.\n", fileName)
		return
	}

	readBytes, readError := templates.ReadFile(strings.ReplaceAll(metadata.RelativeFileNameTemplate, string(filepath.Separator), "/"))
	if nil != readError {
		log.Fatal(readError)
	}

	// Package.
	oldName := []byte("package api")
	newName := []byte("package " + metadata.BaseDirectoryNameCamel)
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// Api path.
	oldName = []byte("/path")
	newName = []byte("/api/" + strings.ReplaceAll(strings.TrimSuffix(metadata.RelativeFileNameCamel, ".go"), string(filepath.Separator), "/"))
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// Api name.
	oldName = []byte("apiName")
	newName = []byte(strings.TrimSuffix(metadata.BaseFileNameTitle, ".go"))
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// Api name.
	oldName = []byte("apiGet")
	newName = []byte(strings.TrimSuffix(metadata.BaseFileNameCamel, ".go") + "Get")
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	writeError := os.WriteFile(fileName, readBytes, os.ModePerm)
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

	fileName := metadata.FullFileName

	if !fileExists(fileName) {
		readBytes, readError := templates.ReadFile(strings.ReplaceAll(metadata.RelativeFileNameTemplate, string(filepath.Separator), "/"))
		if nil != readError {
			log.Fatal(readError)
		}

		// Package.
		oldName := []byte("package pages")
		newName := []byte("package " + metadata.BaseDirectoryNameCamel)
		readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

		// Page name.
		oldName = []byte("pageName")
		newName = []byte(strings.TrimSuffix(metadata.BaseFileNameTitle, ".go"))
		readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

		// Page base.
		oldName = []byte("pageBase")
		newName = []byte(strings.TrimSuffix(metadata.BaseFileNameCamel, ".go") + "Base")
		readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

		// Page action.
		oldName = []byte("pageAction")
		newName = []byte(strings.TrimSuffix(metadata.BaseFileNameCamel, ".go") + "Action")
		readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

		// Page path.
		oldName = []byte("/path")
		newName = []byte("/" + strings.ReplaceAll(strings.TrimSuffix(metadata.RelativeFileNameCamel, ".go"), string(filepath.Separator), "/"))
		readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

		writeError := os.WriteFile(fileName, readBytes, os.ModePerm)
		if writeError != nil {
			log.Fatal(writeError)
		}
	} else {
		fmt.Printf("file `%s` already exists.\n", fileName)
	}

	createViewComponent(metadata.Name)
}

func createViewComponent(pageName string) {
	metadata := findNameMetadataForView(pageName)

	if !fileExists(metadata.FullDirectoryNameCamel) {
		mkdirError := os.MkdirAll(metadata.FullDirectoryNameCamel, os.ModePerm)
		if mkdirError != nil {
			log.Fatal(mkdirError)
		}
	}

	fileName := metadata.FullFileName

	if !fileExists(fileName) {
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
		fmt.Printf("file `%s` already exists.\n", fileName)
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
