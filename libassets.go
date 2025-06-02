package frizzante

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed utilities/*/**
var utilitiesFs embed.FS

//go:embed router/*
var routerFs embed.FS

type AssetsManager interface {
	CreateUtilities() error
	CreateRouter() error
}

type Assets struct {
	utilitiesFs embed.FS
	routerFs    embed.FS
	views       []string
	viewsBase   string
}

func NewAssets() *Assets {
	return &Assets{
		utilitiesFs: utilitiesFs,
		routerFs:    routerFs,
	}
}

func (assets *Assets) LoadViews(from string) error {
	var load func(from string) error

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

func (assets *Assets) CreateUtilities(to string) error {
	var create func(from string, to string) error

	create = func(from string, to string) error {
		embeddedFiles, readDirError := assets.utilitiesFs.ReadDir(from)
		if readDirError != nil {
			return readDirError
		}

		if !IsDirectory(to) {
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

			embeddedContents, readError := utilitiesFs.ReadFile(fromFileName)
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

	return create("utilities", to)
}

func (assets *Assets) CreateRouter(to string) error {
	var create func(from string, to string) error

	create = func(from string, to string) error {
		embeddedFiles, readDirError := assets.routerFs.ReadDir(from)
		if readDirError != nil {
			return readDirError
		}

		if !IsDirectory(to) {
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

			embeddedBytes, readError := routerFs.ReadFile(fromFileName)
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

	return create("router", to)
}
