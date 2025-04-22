package frizzante

import (
	"bufio"
	"bytes"
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*/**
var templates embed.FS

type NameMetadata struct {
	RelativeFileNameTemplate string
	FullDirectoryNameCamel   string
	FullFileNameCamel        string
	RelativeFileNameCamel    string
	FunctionNamePascal       string
	PackageNameCamel         string
}

func findNameMetadata(root string, template string, name string, message string) *NameMetadata {
	if !Exists(root) {
		writeError := os.MkdirAll(root, os.ModePerm)
		if writeError != nil {
			panic(writeError)
		}
	}

	if "" == name {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(message)
		name, _ = reader.ReadString('\n')
		if "" == name {
			return findNameMetadata(root, template, name, message)
		}
	}

	fullRoot, fullFileNameRootError := filepath.Abs(root)
	if nil != fullFileNameRootError {
		panic(fullFileNameRootError)
	}

	metadata := &NameMetadata{}
	trimmedName := strings.Trim(name, "\r\n\t ")
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

	relativeFileNamePascal := ""
	for _, pageNameChunked := range strings.Split(baseFileName, "_") {
		for _, section := range strings.Split(pageNameChunked, string(filepath.Separator)) {
			relativeFileNamePascal = relativeFileNamePascal + string(filepath.Separator) + strings.ToTitle(section[0:1]) + section[1:]
		}
	}
	relativeFileNamePascal = relativeFileNamePascal + extensionName

	relativeFileNameCamel := ""
	for _, pageNameChunked := range strings.Split(baseFileName, "_") {
		for _, section := range strings.Split(pageNameChunked, string(filepath.Separator)) {
			relativeFileNameCamel = relativeFileNameCamel + string(filepath.Separator) + strings.ToLower(section[0:1]) + section[1:]
		}
	}
	relativeFileNameCamel = relativeFileNameCamel[1:] + extensionName

	fullFileNameCamel := filepath.Join(fullRoot, relativeFileNameCamel)

	metadata.FullDirectoryNameCamel = filepath.Dir(fullFileNameCamel)
	metadata.FullFileNameCamel = fullFileNameCamel
	metadata.RelativeFileNameCamel = relativeFileNameCamel
	metadata.RelativeFileNameTemplate = template
	metadata.FunctionNamePascal = filepath.Base(strings.TrimSuffix(relativeFileNamePascal, extensionName))
	metadata.PackageNameCamel = filepath.Base(metadata.FullDirectoryNameCamel)
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

func findNameMetadataForGuard(name string) *NameMetadata {
	return findNameMetadata(
		filepath.Join("lib", "guards"),
		filepath.Join("templates", "guards", "example.go"),
		name,
		"Name the guard: ",
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

func createApi(apiName string) {
	metadata := findNameMetadataForApi(apiName)

	if !Exists(metadata.FullDirectoryNameCamel) {
		mkdirError := os.MkdirAll(metadata.FullDirectoryNameCamel, os.ModePerm)
		if mkdirError != nil {
			panic(mkdirError)
		}
	}

	if Exists(metadata.FullFileNameCamel) {
		fmt.Printf("Api `%s` already exists.\n", metadata.FunctionNamePascal)
		return
	}

	readBytes, readError := templates.ReadFile(strings.ReplaceAll(metadata.RelativeFileNameTemplate, string(filepath.Separator), "/"))
	if nil != readError {
		panic(readError)
	}

	// Package.
	oldName := []byte("package api")
	newName := []byte("package " + metadata.PackageNameCamel)
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// Pattern.
	oldName = []byte("\"GET /\"")
	newName = []byte("\"GET /api/" + strings.ReplaceAll(strings.TrimSuffix(metadata.RelativeFileNameCamel, ".go"), string(filepath.Separator), "/") + "\"")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	// ApiFunction.
	oldName = []byte("func api(")
	newName = []byte("func " + metadata.FunctionNamePascal + "(")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	writeError := os.WriteFile(metadata.FullFileNameCamel, readBytes, os.ModePerm)
	if writeError != nil {
		panic(writeError)
	}
}

func createGuard(guardName string) {
	metadata := findNameMetadataForGuard(guardName)

	if !Exists(metadata.FullDirectoryNameCamel) {
		mkdirError := os.MkdirAll(metadata.FullDirectoryNameCamel, os.ModePerm)
		if mkdirError != nil {
			panic(mkdirError)
		}
	}

	if Exists(metadata.FullFileNameCamel) {
		fmt.Printf("Guard `%s` already exists.\n", metadata.FullFileNameCamel)
		return
	}

	readBytes, readError := templates.ReadFile(strings.ReplaceAll(metadata.RelativeFileNameTemplate, string(filepath.Separator), "/"))
	if nil != readError {
		panic(readError)
	}

	// Package.
	oldName := []byte("package guards")
	newName := []byte("package " + metadata.PackageNameCamel)
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// GuardFunction.
	oldName = []byte("func guard(")
	newName = []byte("func " + metadata.FunctionNamePascal + "(")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	writeError := os.WriteFile(metadata.FullFileNameCamel, readBytes, os.ModePerm)
	if writeError != nil {
		panic(writeError)
	}
}

func createPage(pageName string) {
	metadata := findNameMetadataForPage(pageName)

	if !Exists(metadata.FullDirectoryNameCamel) {
		mkdirError := os.MkdirAll(metadata.FullDirectoryNameCamel, os.ModePerm)
		if mkdirError != nil {
			panic(mkdirError)
		}
	}

	if Exists(metadata.FullFileNameCamel) {
		fmt.Printf("Page function `%s` already exists.\n", metadata.FunctionNamePascal)
		return
	}

	readBytes, readError := templates.ReadFile(strings.ReplaceAll(metadata.RelativeFileNameTemplate, string(filepath.Separator), "/"))
	if nil != readError {
		panic(readError)
	}

	// Package.
	oldName := []byte("package pages")
	newName := []byte("package " + metadata.PackageNameCamel)
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// PageFunction.
	oldName = []byte("func page(")
	newName = []byte("func " + metadata.FunctionNamePascal + "(")
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// Document.
	oldName = []byte("\"pageName\"")
	newName = []byte("\"" + pageName + "\"")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	// Path.
	oldName = []byte("\"/path\"")
	newName = []byte("\"/" + strings.ReplaceAll(strings.TrimSuffix(metadata.RelativeFileNameCamel, ".go"), string(filepath.Separator), "/") + "\"")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	writeError := os.WriteFile(metadata.FullFileNameCamel, readBytes, os.ModePerm)
	if writeError != nil {
		panic(writeError)
	}

	metadata.FullFileNameCamel = strings.TrimSuffix(metadata.FullFileNameCamel, ".go") + ".svelte"
	metadata.RelativeFileNameTemplate = strings.TrimSuffix(metadata.RelativeFileNameTemplate, ".go") + ".svelte"
	metadata.RelativeFileNameCamel = strings.TrimSuffix(metadata.RelativeFileNameCamel, ".go") + ".svelte"
	createPageComponent(metadata)
}

func createPageComponent(metadata *NameMetadata) {
	if Exists(metadata.FullFileNameCamel) {
		fmt.Printf("Page component `%s` already exists.\n", metadata.FunctionNamePascal)
		return
	}

	readBytes, readError := templates.ReadFile(strings.ReplaceAll(metadata.RelativeFileNameTemplate, string(filepath.Separator), "/"))
	if nil != readError {
		panic(readError)
	}

	writeError := os.WriteFile(metadata.FullFileNameCamel, readBytes, os.ModePerm)
	if writeError != nil {
		panic(writeError)
	}
}

// Make makes things.
func Make() {
	api := flag.Bool("api", false, "")
	guard := flag.Bool("guard", false, "")
	page := flag.Bool("page", false, "")
	name := flag.String("name", "", "")
	flag.Parse()

	if *api {
		createApi(*name)
	}

	if *guard {
		createGuard(*name)
	}

	if *page {
		createPage(*name)
	}
}
