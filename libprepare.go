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

// Prepare prepares the `.frizzante` directory.
func Prepare() {
	// Prepare lib.
	err := prepareLib()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare view.
	err = prepareViews()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare vite-project/render.server.svelte.
	err = prepareSsr()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare vite-project/render.client.svelte.
	err = prepareCsr()
	if err != nil {
		log.Fatal(err)
	}
}

func prepareLib() error {
	asyncSvelte, asyncSvelteError := viteProject.ReadFile("vite-project/view.async.svelte")
	if asyncSvelteError != nil {
		return asyncSvelteError
	}

	indexHtml, indexHtmlError := viteProject.ReadFile("vite-project/index.html")
	if indexHtmlError != nil {
		return indexHtmlError
	}

	renderClientJs, renderClientJsError := viteProject.ReadFile("vite-project/render.client.js")
	if renderClientJsError != nil {
		return renderClientJsError
	}

	renderServerJs, renderServerJsError := viteProject.ReadFile("vite-project/render.server.js")
	if renderServerJsError != nil {
		return renderServerJsError
	}

	formSvelte, formSvelteError := viteProject.ReadFile("vite-project/lib/components/Form.svelte")
	if formSvelteError != nil {
		return formSvelteError
	}

	submitSvelte, submitSvelteError := viteProject.ReadFile("vite-project/lib/components/Submit.svelte")
	if submitSvelteError != nil {
		return submitSvelteError
	}

	linkSvelte, linkSvelteError := viteProject.ReadFile("vite-project/lib/components/Link.svelte")
	if linkSvelteError != nil {
		return linkSvelteError
	}

	updateJs, updateJsError := viteProject.ReadFile("vite-project/lib/scripts/update.js")
	if updateJsError != nil {
		return updateJsError
	}

	uuidJs, uuidJsError := viteProject.ReadFile("vite-project/lib/scripts/uuid.js")
	if uuidJsError != nil {
		return uuidJsError
	}

	if !fileExists(".frizzante/vite-project") {
		err := os.MkdirAll(".frizzante/vite-project", os.ModePerm)
		if err != nil {
			return err
		}
	}

	err := os.WriteFile(".frizzante/vite-project/view.async.svelte", asyncSvelte, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(".frizzante/vite-project/index.html", indexHtml, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(".frizzante/vite-project/render.client.js", renderClientJs, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(".frizzante/vite-project/render.server.js", renderServerJs, os.ModePerm)
	if err != nil {
		return err
	}

	if !fileExists(".frizzante/vite-project/lib/components") {
		err = os.MkdirAll(".frizzante/vite-project/lib/components", os.ModePerm)
		if err != nil {
			return err
		}
	}

	if !fileExists(".frizzante/vite-project/lib/scripts") {
		err = os.MkdirAll(".frizzante/vite-project/lib/scripts", os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(".frizzante/vite-project/lib/components/Form.svelte", formSvelte, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(".frizzante/vite-project/lib/components/Submit.svelte", submitSvelte, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(".frizzante/vite-project/lib/components/Link.svelte", linkSvelte, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(".frizzante/vite-project/lib/scripts/update.js", updateJs, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(".frizzante/vite-project/lib/scripts/uuid.js", uuidJs, os.ModePerm)
	if err != nil {
		return err
	}

	if !fileExists(".frizzante/vite-project") {
		err = os.MkdirAll(".frizzante/vite-project", os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func prepareViews() error {
	libViews := filepath.Join("lib", "components", "views")
	sep := string(filepath.Separator)
	suffix := ".svelte"
	return filepath.Walk(
		libViews,
		func(fileName string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || !strings.HasSuffix(info.Name(), suffix) {
				return nil
			}

			fileNameBase := strings.Trim(strings.TrimPrefix(fileName, libViews), sep)
			view := strings.TrimSuffix(strings.ReplaceAll(fileNameBase, sep, "."), ".svelte")

			importFileName, err := filepath.Rel(".frizzante/vite-project", fileName)
			if err != nil {
				log.Fatal(err)
			}
			components[view] = fmt.Sprintf("./%s", importFileName)

			return nil
		},
	)
}

func prepareSsr() error {
	var builder strings.Builder
	renderServerSvelte, readError := viteProject.ReadFile("vite-project/render.server.svelte")
	if readError != nil {
		return readError
	}
	for component, fileName := range components {
		viewAsComponentName := strings.ToUpper(strings.ReplaceAll(component, ".", "_"))
		builder.WriteString(fmt.Sprintf("    import %s from '%s'\n", viewAsComponentName, fileName))
	}

	renderServerSvelteString := strings.Replace(string(renderServerSvelte), "//:app-imports", builder.String(), 1)

	builder.Reset()
	counter := 0
	for view, _ := range components {
		viewAsComponentName := strings.ReplaceAll(view, ".", "_")
		if 0 == counter {
			builder.WriteString(fmt.Sprintf("{#if '%s' === view}\n", view))
		} else {
			builder.WriteString(fmt.Sprintf("{:else if '%s' === view}\n", view))
		}
		builder.WriteString(fmt.Sprintf("    <%s />\n", strings.ToUpper(viewAsComponentName)))
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

func prepareCsr() error {
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
	for view, fileName := range components {
		if 0 == counter {
			builder.WriteString(fmt.Sprintf("{#if '%s' === viewState}\n", view))
		} else {
			builder.WriteString(fmt.Sprintf("{:else if '%s' === viewState}\n", view))
		}
		builder.WriteString(fmt.Sprintf("    <View from={import('%s')} />\n", fileName))
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
