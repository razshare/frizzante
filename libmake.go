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

func createApi(apiName string) {
	fileName := filepath.Join("lib", "components", "server")
	if !Exists(fileName) {
		writeError := os.MkdirAll(fileName, os.ModePerm)
		if writeError != nil {
			panic(writeError)
		}
	}

	if "" == apiName {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Name the api: ")
		apiName, _ = reader.ReadString('\n')
		if "" == apiName {
			createApi(apiName)
			return
		}
	}

	apiName = strings.Trim(strings.ReplaceAll(apiName, "-", "_"), "\r\n\t ")

	apiNameCamel := strings.ToLower(apiName[0:1]) + apiName[1:]
	//apiNamePascal := strings.ToTitle(apiName[0:1]) + apiName[1:]
	newDirectoryName := filepath.Join("lib", "components", "server", apiNameCamel)

	if !Exists(newDirectoryName) {
		mkdirError := os.MkdirAll(newDirectoryName, os.ModePerm)
		if mkdirError != nil {
			panic(mkdirError)
		}
	}

	oldFileName := "templates/api/example.go"
	newFileName := filepath.Join(newDirectoryName, "api.go")
	readBytes, readError := templates.ReadFile(oldFileName)
	if nil != readError {
		panic(readError)
	}
	if Exists(newFileName) {
		fmt.Printf("Api `%s` already exists.\n", apiNameCamel)
		return
	}

	// Package.
	oldName := []byte("package api")
	newName := []byte("package " + apiNameCamel)
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	// Pattern.
	oldName = []byte("\"GET /\"")
	newName = []byte("\"GET /api/" + apiNameCamel + "\"")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	// Api.
	oldName = []byte("func api(")
	newName = []byte("func Api(")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	writeError := os.WriteFile(newFileName, readBytes, os.ModePerm)
	if writeError != nil {
		panic(writeError)
	}

}

func createIndex(indexName string) {
	fileName := filepath.Join("lib", "components", "server")
	if !Exists(fileName) {
		writeError := os.MkdirAll(fileName, os.ModePerm)
		if writeError != nil {
			panic(writeError)
		}
	}

	if "" == indexName {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Name the index: ")
		indexName, _ = reader.ReadString('\n')
		if "" == indexName {
			createIndex(indexName)
			return
		}
	}

	indexName = strings.Trim(strings.ReplaceAll(indexName, "-", "_"), "\r\n\t ")

	indexNameCamel := strings.ToLower(indexName[0:1]) + indexName[1:]
	//indexNamePascal := strings.ToTitle(indexName[0:1]) + indexName[1:]
	newDirectoryName := filepath.Join("lib", "components", "server", indexNameCamel)

	if !Exists(newDirectoryName) {
		mkdirError := os.MkdirAll(newDirectoryName, os.ModePerm)
		if mkdirError != nil {
			panic(mkdirError)
		}
	}

	oldFileName := "templates/indexes/example.go"
	newFileName := filepath.Join(newDirectoryName, "index.go")
	readBytes, readError := templates.ReadFile(oldFileName)
	if nil != readError {
		panic(readError)
	}
	if Exists(newFileName) {
		fmt.Printf("Index `%s` already exists.\n", indexNameCamel)
		return
	}

	// Package.
	oldName := []byte("package indexes")
	newName := []byte("package " + indexNameCamel)
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	// Index.
	oldName = []byte("func index(")
	newName = []byte("func Index(")
	readBytes = bytes.ReplaceAll(readBytes, oldName, newName)

	// Page.
	oldName = []byte("\"page\"")
	newName = []byte("\"" + indexName + "\"")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	// Path.
	oldName = []byte("\"/path\"")
	newName = []byte("\"/" + indexName + "\"")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	writeError := os.WriteFile(newFileName, readBytes, os.ModePerm)
	if writeError != nil {
		panic(writeError)
	}

}

func createGuard(guardName string) {
	fileName := filepath.Join("lib", "components", "server")
	if !Exists(fileName) {
		writeError := os.MkdirAll(fileName, os.ModePerm)
		if writeError != nil {
			panic(writeError)
		}
	}

	if "" == guardName {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Name the guard: ")
		guardName, _ = reader.ReadString('\n')
		if "" == guardName {
			createGuard(guardName)
			return
		}
	}

	guardName = strings.Trim(strings.ReplaceAll(guardName, "-", "_"), "\r\n\t ")

	guardNameCamel := strings.ToLower(guardName[0:1]) + guardName[1:]
	//guardNamePascal := strings.ToTitle(guardName[0:1]) + guardName[1:]
	newDirectoryName := filepath.Join("lib", "components", "server", guardNameCamel)

	if !Exists(newDirectoryName) {
		mkdirError := os.MkdirAll(newDirectoryName, os.ModePerm)
		if mkdirError != nil {
			panic(mkdirError)
		}
	}

	oldFileName := "templates/guards/example.go"
	newFileName := filepath.Join(newDirectoryName, "guard.go")
	readBytes, readError := templates.ReadFile(oldFileName)
	if nil != readError {
		panic(readError)
	}

	if Exists(newFileName) {
		fmt.Printf("Guard `%s` already exists.\n", guardNameCamel)
		return
	}

	// Package.
	oldName := []byte("package guards")
	newName := []byte("package " + guardNameCamel)
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	// Guard.
	oldName = []byte("func guard(")
	newName = []byte("func Guard(")
	readBytes = bytes.Replace(readBytes, oldName, newName, 1)

	writeError := os.WriteFile(newFileName, readBytes, os.ModePerm)
	if writeError != nil {
		panic(writeError)
	}
}

func createPage(pageName string) {
	fileName := filepath.Join("lib", "pages")
	if !Exists(fileName) {
		writeError := os.MkdirAll(fileName, os.ModePerm)
		if writeError != nil {
			panic(writeError)
		}
	}

	if "" == pageName {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Name the page: ")
		pageName, _ = reader.ReadString('\n')
		if "" == pageName {
			createPage(pageName)
			return
		}
	}

	pageName = strings.Trim(strings.ReplaceAll(pageName, "-", "_"), "\r\n\t ")

	//pageNameCamel := strings.ToLower(pageName[0:1])+pageName[1:]
	//pageNamePascal := strings.ToTitle(pageName[0:1]) + pageName[1:]

	oldFileName := "templates/pages/example.svelte"
	newFileName := filepath.Join("lib", "pages", pageName+".svelte")
	readBytes, readError := templates.ReadFile(oldFileName)
	if nil != readError {
		panic(readError)
	}

	if Exists(newFileName) {
		fmt.Printf("Page `%s` already exists.\n", pageName)
		return
	}

	writeError := os.WriteFile(newFileName, readBytes, os.ModePerm)
	if writeError != nil {
		panic(writeError)
	}
}

// Make makes things.
func Make() {
	api := flag.Bool("api", false, "")
	index := flag.Bool("index", false, "")
	guard := flag.Bool("guard", false, "")
	page := flag.Bool("page", false, "")
	name := flag.String("name", "", "")
	flag.Parse()

	if *api {
		createApi(*name)
	}

	if *index {
		createIndex(*name)
	}

	if *guard {
		createGuard(*name)
	}

	if *page {
		createPage(*name)
	}
}
