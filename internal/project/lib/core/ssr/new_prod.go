//go:build prod

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
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

func New(limit int64) renders.Render {
	var mut sync.Mutex
	var server = filepath.Join("app", "dist", "server", "app.server.js")
	var index = filepath.Join("app", "dist", "client", "index.html")
	var jsRenders = make(chan javascript.Render, 1)
	server = strings.ReplaceAll(server, "\\", "/")
	index = strings.ReplaceAll(index, "\\", "/")
	var goRender = func(options renders.RenderOptions) (jsRender javascript.Render, err error) {
		if !embeds.IsFile(options.Efs, server) {
			err = fmt.Errorf("file %s not found", server)
			return
		}
		var data []byte
		if data, err = options.Efs.ReadFile(server); err != nil {
			return
		}
		source := string(data)
		jsRender, err = javascript.NewRender(javascript.NewRenderOptions{
			Server:     server,
			InfoLog:    options.InfoLog,
			ErrorLog:   options.ErrorLog,
			FindSource: func() (string, error) { return source, nil },
		})
		return
	}
	return func(options renders.RenderOptions) (document string, err error) {
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
			var jsRender javascript.Render
			if limit >= 0 {
				mut.Lock()
				if limit >= 0 {
					limit--
				}
				mut.Unlock()
				if jsRender, err = goRender(options); err != nil {
					mut.Lock()
					limit++
					mut.Unlock()
					return
				}
				defer func() { go func() { jsRenders <- jsRender }() }()
			} else {
				jsRender = <-jsRenders
				defer func() { go func() { jsRenders <- jsRender }() }()
			}
			var head string
			var body string
			if head, body, err = jsRender(javascript.RenderOptions{View: view}); err != nil {
				return
			}
			if view.RenderMode == views.RenderModeServer {
				document = renders.NoScript.ReplaceAllString(document, "")
			}
			if view.RenderMode == views.RenderModeServer {
				document = strings.Replace(document, "<!--app-data-->", "", 1)
			} else {
				var data []byte
				if data, err = json.Marshal(options.Data); err != nil {
					return
				}
				document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(renders.DataFormat, data), 1)
			}
			document = strings.Replace(document, "<!--app-head-->", head, 1)
			document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(renders.BodyFormat, body), 1)
			return
		}

		if view.RenderMode == views.RenderModeClient {
			var data []byte
			if data, err = json.Marshal(options.Data); err != nil {
				return
			}
			document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(renders.BodyFormat, ""), 1)
			document = strings.Replace(document, "<!--app-head-->", fmt.Sprintf(renders.HeadFormat, view.Title), 1)
			document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(renders.DataFormat, data), 1)
			return
		}
		err = errors.New("unknown render mode")
		return
	}
}
