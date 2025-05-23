package frizzante

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed vite-project/*
var viteProject embed.FS

func findLibrary(directoryName string) (map[string][]byte, error) {
	contents := map[string][]byte{}
	entries, readDirError := viteProject.ReadDir(directoryName)
	if readDirError != nil {
		return nil, readDirError
	}

	for _, entry := range entries {
		fileName := filepath.Join(directoryName, entry.Name())
		if entry.IsDir() {
			innerContents, findError := findLibrary(fileName)
			if nil != findError {
				return nil, findError
			}

			for key, value := range innerContents {
				contents[key] = value
			}
			continue
		}
		readBytes, readError := viteProject.ReadFile(fileName)
		if readError != nil {
			return nil, readError
		}
		contents[fileName] = readBytes
	}

	return contents, nil
}

func dumpLibrary(library map[string][]byte) error {
	for relativeFileName, content := range library {
		fileName := filepath.Join(".frizzante", relativeFileName)
		directoryName := filepath.Dir(fileName)

		if !fileExists(directoryName) {
			mkdirError := os.MkdirAll(directoryName, os.ModePerm)
			if mkdirError != nil {
				return mkdirError
			}
		}

		writeError := os.WriteFile(fileName, content, os.ModePerm)
		if writeError != nil {
			return writeError
		}
	}

	return nil
}

func findViews() (map[string]string, error) {
	libViews := filepath.Join(PAGES_ROOT)
	sep := string(filepath.Separator)
	suffix := ".svelte"
	views := map[string]string{}
	walkError := filepath.Walk(
		libViews,
		func(fileName string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || !strings.HasSuffix(info.Name(), suffix) {
				return nil
			}

			fileNameRelative := strings.Trim(strings.TrimPrefix(fileName, libViews), sep)
			id := strings.TrimSuffix(strings.TrimSuffix(strings.ReplaceAll(fileNameRelative, sep, "/"), "view.svelte"), "/")
			importFileName, err := filepath.Rel(".frizzante/vite-project", fileName)
			if err != nil {
				return err
			}

			views[id] = fmt.Sprintf("./%s", importFileName)
			return nil
		},
	)

	if walkError != nil {
		return nil, walkError
	}

	return views, nil
}

func dumpSsr(views map[string]string) error {
	var builder strings.Builder
	renderServerSvelte, readError := viteProject.ReadFile("vite-project/render.server.svelte")
	if readError != nil {
		return readError
	}

	for view, fileName := range views {
		componentName := strings.ToUpper(strings.ReplaceAll(view, ".", "_"))
		builder.WriteString(fmt.Sprintf("    import %s from './%s'\n", componentName, fileName))
	}

	renderServerSvelteString := strings.Replace(string(renderServerSvelte), "//:app-imports", builder.String(), 1)

	builder.Reset()
	counter := 0
	for view := range views {
		componentName := strings.ReplaceAll(view, ".", "_")
		if counter == 0 {
			builder.WriteString(fmt.Sprintf("{#if '%s' === server.id}\n", view))
		} else {
			builder.WriteString(fmt.Sprintf("{:else if '%s' === server.id}\n", view))
		}
		builder.WriteString(fmt.Sprintf("    <%s />\n", strings.ToUpper(componentName)))
		counter++
	}
	if counter > 0 {
		builder.WriteString("{/if}")
	}

	renderServerSvelteString = strings.Replace(renderServerSvelteString, "<!--app-router-->", builder.String(), 1)

	writeError := os.WriteFile(".frizzante/vite-project/render.server.svelte", []byte(renderServerSvelteString), os.ModePerm)
	if writeError != nil {
		return writeError
	}
	return nil
}

func dumpCsr(views map[string]string) error {
	// Build client loader.
	renderClientSvelte, readError := viteProject.ReadFile("vite-project/render.client.svelte")
	if readError != nil {
		return readError
	}

	var builder strings.Builder
	builder.WriteString("import View from './view.async.svelte'")
	renderClientSvelteString := strings.Replace(string(renderClientSvelte), "//:app-imports", builder.String(), 1)

	builder.Reset()
	counter := 0
	for view, fileName := range views {
		if counter == 0 {
			builder.WriteString(fmt.Sprintf("{#if '%s' === server.id}\n", view))
		} else {
			builder.WriteString(fmt.Sprintf("{:else if '%s' === server.id}\n", view))
		}
		builder.WriteString(fmt.Sprintf("    <View from={import('./%s')} />\n", fileName))
		counter++
	}
	if counter > 0 {
		builder.WriteString("{/if}")
	}
	renderClientSvelteString = strings.Replace(renderClientSvelteString, "<!--app-router-->", builder.String(), 1)

	// Dump client loader.
	writeError := os.WriteFile(".frizzante/vite-project/render.client.svelte", []byte(renderClientSvelteString), os.ModePerm)
	if writeError != nil {
		return writeError
	}

	return nil
}

// Prepare prepares the `.frizzante` directory.
func Prepare() {
	// Find library.
	library, libraryError := findLibrary("vite-project")
	if nil != libraryError {
		log.Fatal(libraryError)
	}

	// Dump library.
	err := dumpLibrary(library)
	if err != nil {
		log.Fatal(err)
	}

	// Find views.
	views, viewsError := findViews()
	if nil != viewsError {
		log.Fatal(viewsError)
	}

	// Dump vite-project/render.server.svelte.
	err = dumpSsr(views)
	if err != nil {
		log.Fatal(err)
	}

	// Dump vite-project/render.client.svelte.
	err = dumpCsr(views)
	if err != nil {
		log.Fatal(err)
	}
}
