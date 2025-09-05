package ssr

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"main/lib/core/embeds"
	"main/lib/core/js"
	_view "main/lib/core/view"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
)

//go:embed render.format
var RenderFormat string

//go:embed target.format
var TargetFormat string

//go:embed head.format
var HeadFormat string

//go:embed body.format
var BodyFormat string

//go:embed data.format
var DataFormat string

var NoScript = regexp.MustCompile(`<script.*>.*</script>`)

func New(conf Config) func(view _view.View) (html string, err error) {
	var efs = conf.Efs
	var app = conf.App
	var disk = conf.Disk
	var limit = conf.Limit

	if limit <= 0 {
		limit = 1
	}

	if app == "" {
		app = "app"
	}

	var mut sync.Mutex
	var id = "svelte-app"
	var dist = filepath.Join(app, "dist")
	var appServer = filepath.Join(dist, "app.server.js")
	var appServerFix = strings.ReplaceAll(appServer, "\\", "/")
	var index = filepath.Join(dist, "client", "index.html")
	var indexFix = strings.ReplaceAll(index, "\\", "/")
	var renders = make(chan goja.Callable, 1)
	var runtimes = make(chan *goja.Runtime, 1)
	var compile = func() (render goja.Callable, runtime *goja.Runtime, err error) {
		var data []byte

		if !disk && embeds.IsFile(efs, appServerFix) {
			data, err = efs.ReadFile(appServerFix)
		} else {
			data, err = os.ReadFile(appServer)
		}

		if err != nil {
			return
		}

		runtime = goja.New()

		var text string
		if text, err = js.Bundle(app, api.FormatCommonJS, string(data)); err != nil {
			return
		}

		var prog *goja.Program
		if prog, err = goja.Compile(appServer, fmt.Sprintf(RenderFormat, text), false); err != nil {
			return
		}

		var value goja.Value
		if value, err = runtime.RunProgram(prog); err != nil {
			return
		}

		var isfun bool
		if render, isfun = goja.AssertFunction(value); !isfun {
			err = errors.New("render is not a function")
		}

		return
	}

	return func(view _view.View) (html string, err error) {
		var data []byte

		if !disk && embeds.IsFile(efs, indexFix) {
			data, err = efs.ReadFile(indexFix)
		} else {
			data, err = os.ReadFile(index)
		}

		if err != nil {
			return "", err
		}

		html = string(data)

		if view.RenderMode == _view.RenderModeServer || view.RenderMode == _view.RenderModeFull {
			var render goja.Callable
			var runtime *goja.Runtime
			if disk {
				render, runtime, err = compile()
				if err != nil {
					return
				}
			} else if limit >= 0 {
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

			var promise goja.Value
			if promise, err = render(goja.Undefined(), runtime.ToValue(_view.Data(view))); err != nil {
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

			if view.RenderMode == _view.RenderModeServer {
				html = NoScript.ReplaceAllString(html, "")
			}

			if view.RenderMode == _view.RenderModeServer {
				html = strings.Replace(html, "<!--app-target-->", "", 1)
				html = strings.Replace(html, "<!--app-data-->", "", 1)
			} else {
				if data, err = json.Marshal(_view.Data(view)); err != nil {
					return
				}

				html = strings.Replace(html, "<!--app-target-->", fmt.Sprintf(TargetFormat, id), 1)
				html = strings.Replace(html, "<!--app-data-->", fmt.Sprintf(DataFormat, data), 1)
			}

			html = strings.Replace(html, "<!--app-head-->", head, 1)
			html = strings.Replace(html, "<!--app-body-->", fmt.Sprintf(BodyFormat, id, body), 1)

			return
		}

		if view.RenderMode == _view.RenderModeClient {
			if data, err = json.Marshal(_view.Data(view)); err != nil {
				return
			}

			html = strings.Replace(html, "<!--app-target-->", fmt.Sprintf(TargetFormat, id), 1)
			html = strings.Replace(html, "<!--app-body-->", fmt.Sprintf(BodyFormat, id, ""), 1)
			html = strings.Replace(html, "<!--app-head-->", fmt.Sprintf(HeadFormat, view.Title), 1)
			html = strings.Replace(html, "<!--app-data-->", fmt.Sprintf(DataFormat, data), 1)

			return
		}

		err = errors.New("unknown render mode")

		return
	}
}
