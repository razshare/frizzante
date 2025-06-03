package cli

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/razshare/frizzante/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed utilities/*
var utilitiesEfs embed.FS

type AotUtilities struct {
	efs embed.FS
}

func NewAotUtilities() *AotUtilities {
	return &AotUtilities{
		efs: utilitiesEfs,
	}
}

func (assets *AotUtilities) CreateOnDisk(to string) error {
	var from = "utilities"
	//var to = filepath.Join("frz/.generated", "utilities")
	var create func(from string, to string) error

	create = func(from string, to string) error {
		embeddedFiles, readDirError := assets.efs.ReadDir(from)
		if readDirError != nil {
			return readDirError
		}

		if !fs.IsDirectory(to) {
			mkdirError := os.MkdirAll(to, os.ModePerm)
			if mkdirError != nil {
				return mkdirError
			}
		}

		for _, embeddedFile := range embeddedFiles {
			fromFileName := strings.Join([]string{from, embeddedFile.Name()}, "/")
			toFileName := filepath.Join(to, strings.ReplaceAll(embeddedFile.Name(), "/", string(filepath.Separator)))
			if embeddedFile.IsDir() {
				createError := create(fromFileName, toFileName)
				if createError != nil {
					return createError
				}
				continue
			}

			embeddedContents, readError := utilitiesEfs.ReadFile(fromFileName)
			if readError != nil {
				return readError
			}

			writeError := os.WriteFile(toFileName, embeddedContents, os.ModePerm)
			if writeError != nil {
				return writeError
			}
		}
		return nil
	}

	return create(from, to)
}

//go:embed router/*
var routerEfs embed.FS

type AotRouter struct {
	efs       embed.FS
	views     []string
	viewsBase string
}

func NewAotRouter() *AotRouter {
	return &AotRouter{
		efs: routerEfs,
	}
}

func (assets *AotRouter) LoadViews(from string) error {
	//var from = filepath.Join("lib", "components", "views")
	var load func(from string) error
	assets.views = []string{}

	load = func(from string) error {
		views, readDirError := os.ReadDir(from)
		if readDirError != nil {
			return readDirError
		}

		for _, view := range views {
			fileName := filepath.Join(from, view.Name())

			if view.IsDir() {
				loadError := load(fileName)
				if loadError != nil {
					return loadError
				}
				continue
			}

			fileNameRelative, fileNameRelativeError := filepath.Rel(from, fileName)
			if fileNameRelativeError != nil {
				return fileNameRelativeError
			}

			assets.views = append(assets.views, fileNameRelative)
		}

		assets.viewsBase = from

		return nil
	}

	return load(from)
}

func (assets *AotRouter) CreateOnDisk(to string) error {
	var from = "router"
	//var to = filepath.Join("frz/.generated", "router")
	var create func(from string, to string) error

	create = func(from string, to string) error {
		embeddedFiles, readDirError := assets.efs.ReadDir(from)
		if readDirError != nil {
			return readDirError
		}

		if !fs.IsDirectory(to) {
			mkdirError := os.MkdirAll(to, os.ModePerm)
			if mkdirError != nil {
				return mkdirError
			}
		}

		for _, embeddedFile := range embeddedFiles {
			fromFileName := strings.Join([]string{from, embeddedFile.Name()}, "/")
			toFileName := filepath.Join(to, strings.ReplaceAll(embeddedFile.Name(), "/", string(filepath.Separator)))
			if embeddedFile.IsDir() {
				createError := create(fromFileName, toFileName)
				if createError != nil {
					return createError
				}
				continue
			}

			embeddedBytes, readError := routerEfs.ReadFile(fromFileName)
			if readError != nil {
				return readError
			}

			// BEGIN MODIFYING ROUTER
			if filepath.Join(to, "server.svelte") == toFileName {
				var imports strings.Builder
				var router strings.Builder

				counter := 0
				for _, fileName := range assets.views {
					importDirectoryRelative, importDirectoryRelativeError := filepath.Rel(to, assets.viewsBase)
					if importDirectoryRelativeError != nil {
						return importDirectoryRelativeError
					}

					importFileNameRelative := strings.Join([]string{importDirectoryRelative, fileName}, "/")

					viewName := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSuffix(fileName, ".svelte"), "\\", "."), "/", "."), "_.")
					componentName := strings.Trim(strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(viewName, "/", "__"), ".", "_")), "_.")

					if counter == 0 {
						router.WriteString(fmt.Sprintf("{#if '%s' === view.name}", viewName))
					} else {
						imports.WriteString("\n    ")
						router.WriteString(fmt.Sprintf("\n{:else if '%s' === view.name}", viewName))
					}
					imports.WriteString(fmt.Sprintf("import %s from './%s'", componentName, importFileNameRelative))
					router.WriteString(fmt.Sprintf("\n    <%s/>", componentName))
					counter++
				}

				if counter > 0 {
					router.WriteString("\n{/if}")
				}

				embeddedBytes = bytes.ReplaceAll(embeddedBytes, []byte("//:imports"), []byte(imports.String()))
				embeddedBytes = bytes.ReplaceAll(embeddedBytes, []byte("<!--:router-->"), []byte(router.String()))
			} else if filepath.Join(to, "client.svelte") == toFileName {
				var imports strings.Builder
				var router strings.Builder

				imports.WriteString("import ViewComponent from './view.svelte'")

				counter := 0
				for _, fileName := range assets.views {
					importDirectoryRelative, importDirectoryRelativeError := filepath.Rel(to, assets.viewsBase)
					if importDirectoryRelativeError != nil {
						return importDirectoryRelativeError
					}

					importFileNameRelative := strings.Join([]string{importDirectoryRelative, fileName}, "/")

					viewName := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSuffix(fileName, ".svelte"), "\\", "."), "/", "."), "_.")

					if counter == 0 {
						router.WriteString(fmt.Sprintf("{#if '%s' === view.name}", viewName))
					} else {
						router.WriteString(fmt.Sprintf("\n{:else if '%s' === view.name}", viewName))
					}
					router.WriteString(fmt.Sprintf("\n    <ViewComponent from={import('./%s')}/>", importFileNameRelative))
					counter++
				}

				if counter > 0 {
					router.WriteString("\n{/if}")
				}

				embeddedBytes = bytes.ReplaceAll(embeddedBytes, []byte("//:imports"), []byte(imports.String()))
				embeddedBytes = bytes.ReplaceAll(embeddedBytes, []byte("<!--:router-->"), []byte(router.String()))
			}
			// END MODIFYING ROUTER

			writeError := os.WriteFile(toFileName, embeddedBytes, os.ModePerm)
			if writeError != nil {
				return writeError
			}
		}
		return nil
	}

	return create(from, to)
}

func CreateAotRouterOnDisk(to string) {
	router := NewAotRouter()
	routerError := router.CreateOnDisk(to)
	if routerError != nil {
		log.Fatal(routerError)
	}
}

func CreateAotUtilitiesOnDisk(to string) {
	utilities := NewAotUtilities()
	utilitiesError := utilities.CreateOnDisk(to)
	if utilitiesError != nil {
		log.Fatal(utilitiesError)
	}
}
