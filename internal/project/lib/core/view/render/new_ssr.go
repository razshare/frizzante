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

	"github.com/dop251/goja"
	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
	"github.com/razshare/frizzante/internal/project/lib/core/types"
	view_ "github.com/razshare/frizzante/internal/project/lib/core/view"
	compile_ "github.com/razshare/frizzante/internal/project/lib/core/view/render/compile"
)

func New(config Config) func(view view_.View) (html string, err error) {
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
		if limitString := os.Getenv("FRIZZANTE_RENDER_LIMIT"); limitString != "" {
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
	var server = filepath.Join(app, "dist", "app.server.js")
	var index = filepath.Join(app, "dist", "client", "index.html")
	var renders = make(chan goja.Callable, 1)
	var runtimes = make(chan *goja.Runtime, 1)

	server = strings.ReplaceAll(server, "\\", "/")
	index = strings.ReplaceAll(index, "\\", "/")

	var compile = func() (render goja.Callable, runtime *goja.Runtime, err error) {
		if !embeds.IsFile(efs, server) {
			err = fmt.Errorf("file %s not found", server)
			return
		}

		var data []byte
		if data, err = efs.ReadFile(server); err != nil {
			return
		}

		render, runtime, err = compile_.New(compile_.Config{
			Data:     data,
			Format:   RenderFormat,
			App:      app,
			Server:   server,
			ErrorLog: errorLog,
			InfoLog:  infoLog,
		})
		return
	}

	return func(view view_.View) (indexString string, err error) {
		var propsData []byte

		if embeds.IsFile(efs, index) {
			propsData, err = efs.ReadFile(index)
		}

		if err != nil {
			return
		}

		indexString = string(propsData)

		if view.RenderMode == view_.RenderModeServer || view.RenderMode == view_.RenderModeFull {
			var render goja.Callable
			var runtime *goja.Runtime
			if limit >= 0 {
				mut.Lock()
				if limit >= 0 {
					limit--
				}
				mut.Unlock()
				render, runtime, err = compile()
				if err != nil {
					return
				}
				defer func() { go func() { renders <- render }() }()
				defer func() { go func() { runtimes <- runtime }() }()
			} else {
				render = <-renders
				runtime = <-runtimes
				defer func() { go func() { renders <- render }() }()
				defer func() { go func() { runtimes <- runtime }() }()
			}

			var props map[string]any
			if props, err = types.EncodeInterface(view_.NewData(view)); err != nil {
				return
			}

			var promise goja.Value
			if promise, err = render(goja.Undefined(), runtime.ToValue(props)); err != nil {
				return
			}

			result := promise.Export().(*goja.Promise).Result().ToObject(runtime)

			headv := result.Get("head")
			bodyv := result.Get("body")

			var head string
			var body string

			if headv != nil {
				head = headv.String()
			}

			if bodyv != nil {
				body = bodyv.String()
			}

			if view.RenderMode == view_.RenderModeServer {
				indexString = NoScript.ReplaceAllString(indexString, "")
			}

			if view.RenderMode == view_.RenderModeServer {
				indexString = strings.Replace(indexString, "<!--app-data-->", "", 1)
			} else {
				if propsData, err = json.Marshal(view_.NewData(view)); err != nil {
					return
				}

				indexString = strings.Replace(indexString, "<!--app-data-->", fmt.Sprintf(DataFormat, propsData), 1)
			}

			indexString = strings.Replace(indexString, "<!--app-head-->", head, 1)
			indexString = strings.Replace(indexString, "<!--app-body-->", fmt.Sprintf(BodyFormat, body), 1)

			return
		}

		if view.RenderMode == view_.RenderModeClient {
			if propsData, err = json.Marshal(view_.NewData(view)); err != nil {
				return
			}

			indexString = strings.Replace(indexString, "<!--app-body-->", fmt.Sprintf(BodyFormat, ""), 1)
			indexString = strings.Replace(indexString, "<!--app-head-->", fmt.Sprintf(HeadFormat, view.Title), 1)
			indexString = strings.Replace(indexString, "<!--app-data-->", fmt.Sprintf(DataFormat, propsData), 1)

			return
		}

		err = errors.New("unknown render mode")

		return
	}
}
