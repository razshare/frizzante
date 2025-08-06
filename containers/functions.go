package containers

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/javascript"
	"log"
	"os"
	"strings"
	"sync"
)

// NewViewContainer creates a new view container that points to the "app" directory.
func NewViewContainer() *ViewContainer {
	return &ViewContainer{
		AppRoot:               "app",
		ServerJs:              "app/dist/server.js",
		IndexHtml:             "app/dist/client/index.html",
		MaximumProgramCounter: 2,
		MaximumRuntimeCounter: 2,
		ProgramChannel:        make(chan *goja.Program, 1),
		RuntimeChannel:        make(chan *goja.Runtime, 1),
	}
}

// ReadServerJs reads the contents of the server script and returns it.
//
// The server script is usually located at "app/dist/server.js".
func (container *ViewContainer) ReadServerJs() (string, error) {
	var data []byte
	var readError error

	if files.IsFile(container.ServerJs) {
		data, readError = os.ReadFile(container.ServerJs)
		if readError != nil {
			return "", readError
		}
	} else {
		data, readError = container.Efs.ReadFile(strings.ReplaceAll(container.ServerJs, "\\", "/"))
		if readError != nil {
			return "", readError
		}
	}

	bundledSourceCode, bundleError := javascript.Bundle(container.AppRoot, api.FormatCommonJS, string(data))
	if bundleError != nil {
		return "", bundleError
	}

	return fmt.Sprintf(
		`
			if (!module) {
				var module={exports:{}}; 
			}

			(function(){
				%s
				return render;
			})()
		`,
		bundledSourceCode,
	), nil
}

// ReadIndexHtml reads the contents of the index document and returns it.
//
// The index document is usually located at "app/dist/client/index.html".
func (container *ViewContainer) ReadIndexHtml() (string, error) {
	if container.IndexHtmlCache != "" {
		return container.IndexHtmlCache, nil
	}

	if files.IsFile(container.IndexHtml) {
		data, readError := os.ReadFile(container.IndexHtml)
		if readError != nil {
			return "", readError
		}

		container.IndexHtmlCache = string(data)

		return container.IndexHtmlCache, nil
	}

	var data []byte
	fileNameFixed := strings.ReplaceAll(container.IndexHtml, "\\", "/")
	if embeds.IsFile(container.Efs, fileNameFixed) {
		var readError error
		data, readError = container.Efs.ReadFile(fileNameFixed)
		if readError != nil {
			return "", readError
		}
	} else {
		return "", errors.New("view index is missing from the host file system and the embedded file system")
	}

	container.IndexHtmlCache = string(data)
	return container.IndexHtmlCache, nil
}

// Start starts the view container.
//
// This will launch two goroutines that sends values to RuntimeChannel and ProgramChannel.
func (container *ViewContainer) Start() {
	var group sync.WaitGroup
	group.Add(2)

	container.Stop = false

	go func() {
		defer group.Done()
		var count uint64
		for {
			if container.Stop || count >= container.MaximumRuntimeCounter {
				return
			}
			container.RuntimeChannel <- goja.New()
			count++
		}
	}()

	go func() {
		defer group.Done()
		var count uint64
		for {
			if container.Stop || count >= container.MaximumProgramCounter {
				return
			}
			sourceCode, readError := container.ReadServerJs()
			if readError != nil {
				log.Fatal(readError)
			}

			program, compileError := goja.Compile("goja", sourceCode, false)
			if compileError != nil {
				log.Fatal(compileError)
			}

			container.ProgramChannel <- program
			count++
		}
	}()

	group.Wait()
}
