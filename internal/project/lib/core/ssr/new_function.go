//go:build !dev

package ssr

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
	"github.com/razshare/frizzante/internal/project/lib/core/javascript"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	render_ "github.com/razshare/frizzante/internal/project/lib/core/views/render"
)

func NewFunction(limit int64) render_.Function {
	var mut sync.Mutex
	var server = filepath.Join("app", "dist", "app.server.cjs")
	var index = filepath.Join("app", "dist", "client", "index.html")
	var renders = make(chan javascript.RenderFunction, 1)

	server = strings.ReplaceAll(server, "\\", "/")
	index = strings.ReplaceAll(index, "\\", "/")

	var renderDocument = func(options render_.FunctionOptions) (render javascript.RenderFunction, err error) {
		if !embeds.IsFile(options.Efs, server) {
			err = fmt.Errorf("file %s not found", server)
			return
		}

		var data []byte
		if data, err = options.Efs.ReadFile(server); err != nil {
			return
		}

		render, err = javascript.NewRenderFunction(javascript.NewRenderFunctionOptions{
			Data:     data,
			Server:   server,
			InfoLog:  options.InfoLog,
			ErrorLog: options.ErrorLog,
		})
		return
	}

	return func(options render_.FunctionOptions) (document string, err error) {
		if !embeds.IsFile(options.Efs, index) {
			err = fmt.Errorf("file %s not found", index)
			return
		}

		var indexData []byte
		if indexData, err = options.Efs.ReadFile(index); err != nil {
			return
		}

		document = string(indexData)

		view := options.View

		if view.RenderMode == views.RenderModeServer || view.RenderMode == views.RenderModeFull {
			var renderContent javascript.RenderFunction
			if limit >= 0 {
				mut.Lock()
				if limit >= 0 {
					limit--
				}
				mut.Unlock()

				if renderContent, err = renderDocument(options); err != nil {
					mut.Lock()
					limit++
					mut.Unlock()
					return
				}
				defer func() { go func() { renders <- renderContent }() }()
			} else {
				renderContent = <-renders
				defer func() { go func() { renders <- renderContent }() }()
			}

			var head string
			var body string
			if head, body, err = renderContent(javascript.RenderFunctionOptions{View: view}); err != nil {
				return
			}

			if view.RenderMode == views.RenderModeServer {
				document = render_.NoScript.ReplaceAllString(document, "")
			}

			if view.RenderMode == views.RenderModeServer {
				document = strings.Replace(document, "<!--app-data-->", "", 1)
			} else {
				var data []byte
				if data, err = json.Marshal(views.NewData(view)); err != nil {
					return
				}

				document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(render_.DataFormat, data), 1)
			}

			document = strings.Replace(document, "<!--app-head-->", head, 1)
			document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(render_.BodyFormat, body), 1)

			return
		}

		if view.RenderMode == views.RenderModeClient {
			var data []byte
			if data, err = json.Marshal(views.NewData(view)); err != nil {
				return
			}

			document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(render_.BodyFormat, ""), 1)
			document = strings.Replace(document, "<!--app-head-->", fmt.Sprintf(render_.HeadFormat, view.Title), 1)
			document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(render_.DataFormat, data), 1)

			return
		}

		err = errors.New("unknown render mode")

		return
	}
}
