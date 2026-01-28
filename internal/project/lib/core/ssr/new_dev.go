//go:build dev

package ssr

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/javascript"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

func New(_ int64) renders.Render {
	var server = filepath.Join("app", "dist", "server", "app.server.cjs")
	var index = filepath.Join("app", "dist", "client", "index.html")
	server = strings.ReplaceAll(server, "/", string(filepath.Separator))
	server = strings.ReplaceAll(server, "\\", string(filepath.Separator))
	index = strings.ReplaceAll(index, "/", string(filepath.Separator))
	index = strings.ReplaceAll(index, "\\", string(filepath.Separator))
	var goRender = func(options renders.RenderOptions) (jsRender javascript.Render, err error) {
		if !files.IsFile(server) {
			err = fmt.Errorf("file %s not found", server)
			return
		}
		var data []byte
		if data, err = os.ReadFile(server); err != nil {
			return
		}
		jsRender, err = javascript.NewRender(javascript.NewRenderOptions{
			Data:     data,
			Server:   server,
			InfoLog:  options.InfoLog,
			ErrorLog: options.ErrorLog,
		})
		return
	}
	return func(options renders.RenderOptions) (document string, err error) {
		if !files.IsFile(index) {
			err = fmt.Errorf("file %s not found", index)
			return
		}
		var indexData []byte
		if indexData, err = os.ReadFile(index); err != nil {
			return
		}
		document = string(indexData)
		view := options.View
		if view.RenderMode == views.RenderModeServer || view.RenderMode == views.RenderModeFull {
			var render javascript.Render
			if render, err = goRender(options); err != nil {
				return
			}
			var head string
			var body string
			if head, body, err = render(javascript.RenderOptions{View: view}); err != nil {
				return
			}
			if view.RenderMode == views.RenderModeServer {
				document = renders.NoScript.ReplaceAllString(document, "")
			}
			if view.RenderMode == views.RenderModeServer {
				document = strings.Replace(document, "<!--app-data-->", "", 1)
			} else {
				var data []byte
				if data, err = json.Marshal(views.NewData(view)); err != nil {
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
			if data, err = json.Marshal(views.NewData(view)); err != nil {
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
