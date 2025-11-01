//go:build !dev && !no_js_runtime

package render

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/render_function"
)

func New(config Config) Render {
	var efs = config.Efs
	var app = config.App
	var limit = config.Limit
	var errorLog = config.ErrorLog
	var infoLog = config.InfoLog

	if errorLog == nil {
		errorLog = log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime)
	}

	if infoLog == nil {
		infoLog = log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime)
	}

	if app == "" {
		app = "app"
	}

	if limit <= 0 {
		if limitString := os.Getenv("FRIZZANTE_JS_RUNTIME_LIMIT"); limitString != "" {
			var err error
			var limit64 int64
			if limit64, err = strconv.ParseInt(limitString, 10, 64); err != nil {
				errorLog.Printf("could not parse frizzante render limit value %s, falling back to limit 1", limitString)
				limit = 1
			} else {
				limit = int(limit64)
			}
		} else {
			limit = 1
		}
	}

	var mut sync.Mutex
	var server = filepath.Join(app, "dist", "app.server.cjs")
	var index = filepath.Join(app, "dist", "client", "index.html")
	var renders = make(chan render_function.RenderFunction, 1)

	server = strings.ReplaceAll(server, "\\", "/")
	index = strings.ReplaceAll(index, "\\", "/")

	var compile = func() (render render_function.RenderFunction, err error) {
		if !embeds.IsFile(efs, server) {
			err = fmt.Errorf("file %s not found", server)
			return
		}

		var data []byte
		if data, err = efs.ReadFile(server); err != nil {
			return
		}

		render, err = render_function.New(render_function.Config{
			Data:     data,
			App:      app,
			Server:   server,
			ErrorLog: errorLog,
			InfoLog:  infoLog,
		})
		return
	}

	return func(view views.View) (document string, err error) {
		if !embeds.IsFile(efs, index) {
			err = fmt.Errorf("file %s not found", index)
			return
		}

		var indexData []byte
		if indexData, err = efs.ReadFile(index); err != nil {
			return
		}

		document = string(indexData)

		if view.RenderMode == views.RenderModeServer || view.RenderMode == views.RenderModeFull {
			var render render_function.RenderFunction
			if limit >= 0 {
				mut.Lock()
				if limit >= 0 {
					limit--
				}
				mut.Unlock()

				if render, err = compile(); err != nil {
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
	}
}
