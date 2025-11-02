//go:build !dev && !no_js_runtime

package render

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/render_function"
)

func New() (Render, error) {
	var limit int
	var mut sync.Mutex
	var server = filepath.Join("app", "dist", "app.server.cjs")
	var index = filepath.Join("app", "dist", "client", "index.html")
	var renders = make(chan render_function.RenderFunction, 1)

	if limitString := os.Getenv("FRIZZANTE_JS_RUNTIME_LIMIT"); limitString != "" {
		var err error
		var limit64 int64
		if limit64, err = strconv.ParseInt(limitString, 10, 64); err != nil {
			return nil, err
		} else {
			limit = int(limit64)
		}
	} else {
		limit = 1
	}

	server = strings.ReplaceAll(server, "\\", "/")
	index = strings.ReplaceAll(index, "\\", "/")

	var compile = func(options Options) (render render_function.RenderFunction, err error) {
		if !embeds.IsFile(options.Efs, server) {
			err = fmt.Errorf("file %s not found", server)
			return
		}

		var data []byte
		if data, err = options.Efs.ReadFile(server); err != nil {
			return
		}

		render, err = render_function.New(render_function.Options{
			Data:     data,
			Server:   server,
			InfoLog:  options.InfoLog,
			ErrorLog: options.ErrorLog,
		})
		return
	}

	return func(options Options) (document string, err error) {
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
			var render render_function.RenderFunction
			if limit >= 0 {
				mut.Lock()
				if limit >= 0 {
					limit--
				}
				mut.Unlock()

				if render, err = compile(options); err != nil {
					mut.Lock()
					limit++
					mut.Unlock()
					return
				}
				defer func() { go func() { renders <- render }() }()
			} else {
				render = <-renders
				defer func() { go func() { renders <- render }() }()
			}

			var head string
			var body string
			if head, body, err = render(view); err != nil {
				return
			}

			if view.RenderMode == views.RenderModeServer {
				document = NoScript.ReplaceAllString(document, "")
			}

			if view.RenderMode == views.RenderModeServer {
				document = strings.Replace(document, "<!--app-data-->", "", 1)
			} else {
				var data []byte
				if data, err = json.Marshal(views.NewData(view)); err != nil {
					return
				}

				document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(DataFormat, data), 1)
			}

			document = strings.Replace(document, "<!--app-head-->", head, 1)
			document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(BodyFormat, body), 1)

			return
		}

		if view.RenderMode == views.RenderModeClient {
			var data []byte
			if data, err = json.Marshal(views.NewData(view)); err != nil {
				return
			}

			document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(BodyFormat, ""), 1)
			document = strings.Replace(document, "<!--app-head-->", fmt.Sprintf(HeadFormat, view.Title), 1)
			document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(DataFormat, data), 1)

			return
		}

		err = errors.New("unknown render mode")

		return
	}, nil
}
