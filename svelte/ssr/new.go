package ssr

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/js"
	"github.com/razshare/frizzante/view"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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

func New(conf Config) view.Render {
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
	var nameDist = filepath.Join(app, "dist")
	var nameScript = filepath.Join(nameDist, "server.js")
	var nameScriptFixed = strings.ReplaceAll(nameScript, "\\", "/")
	var nameDoc = filepath.Join(nameDist, "client", "index.html")
	var nameDocFixed = strings.ReplaceAll(nameDoc, "\\", "/")
	var renders = make(chan goja.Callable, 1)
	var runtimes = make(chan *goja.Runtime, 1)
	var compile = func() (render goja.Callable, runtime *goja.Runtime, err error) {
		var data []byte

		if !disk && embeds.IsFile(efs, nameScriptFixed) {
			data, err = efs.ReadFile(nameScriptFixed)
		} else {
			data, err = os.ReadFile(nameScript)
		}

		if err != nil {
			return
		}

		runtime = goja.New()

		var source string
		if source, err = js.Bundle(app, api.FormatCommonJS, string(data)); err != nil {
			return
		}

		var prog *goja.Program
		if prog, err = goja.Compile(nameScript, fmt.Sprintf(RenderFormat, source), false); err != nil {
			return
		}

		var v goja.Value
		if v, err = runtime.RunProgram(prog); err != nil {
			return
		}

		var isf bool
		if render, isf = goja.AssertFunction(v); !isf {
			err = errors.New("render is not a function")
		}

		return
	}

	return func(v view.View) (html string, err error) {
		var data []byte

		if !disk && embeds.IsFile(efs, nameDocFixed) {
			data, err = efs.ReadFile(nameDocFixed)
		} else {
			data, err = os.ReadFile(nameDoc)
		}

		if err != nil {
			return "", err
		}

		html = string(data)

		if v.Render == view.RenderServer || v.Render == view.RenderFull {
			var render goja.Callable
			var runtime *goja.Runtime
			if disk {
				render, runtime, err = compile()
			} else if limit >= 0 {
				mut.Lock()
				if limit >= 0 {
					limit--
				}
				mut.Unlock()
				render, runtime, err = compile()
				defer func() { go func() { renders <- render }() }()
				defer func() { go func() { runtimes <- runtime }() }()
			} else {
				render = <-renders
				runtime = <-runtimes
				defer func() { go func() { renders <- render }() }()
				defer func() { go func() { runtimes <- runtime }() }()
			}

			if err != nil {
				return "", err
			}

			promise, perr := render(goja.Undefined(), runtime.ToValue(view.Data(v)))

			if perr != nil {
				return "", perr
			}

			value := promise.Export().(*goja.Promise).Result().ToObject(runtime)

			headv := value.Get("head")
			bodyv := value.Get("body")

			var head string
			var body string

			if headv != nil {
				head = headv.String()
			}

			if bodyv != nil {
				body = bodyv.String()
			}

			if v.Render == view.RenderServer {
				html = NoScript.ReplaceAllString(html, "")
			}

			if v.Render == view.RenderServer {
				html = strings.Replace(html, "<!--app-target-->", "", 1)
				html = strings.Replace(html, "<!--app-data-->", "", 1)
			} else {
				if data, err = json.Marshal(view.Data(v)); err != nil {
					return
				}

				html = strings.Replace(html, "<!--app-target-->", fmt.Sprintf(TargetFormat, id), 1)
				html = strings.Replace(html, "<!--app-data-->", fmt.Sprintf(DataFormat, data), 1)
			}

			html = strings.Replace(html, "<!--app-head-->", head, 1)
			html = strings.Replace(html, "<!--app-body-->", fmt.Sprintf(BodyFormat, id, body), 1)

			return html, nil
		}

		if v.Render == view.RenderClient {
			if data, err = json.Marshal(view.Data(v)); err != nil {
				return
			}

			html = strings.Replace(html, "<!--app-target-->", fmt.Sprintf(TargetFormat, id), 1)
			html = strings.Replace(html, "<!--app-body-->", fmt.Sprintf(BodyFormat, id, ""), 1)
			html = strings.Replace(html, "<!--app-head-->", fmt.Sprintf(HeadFormat, v.Title), 1)
			html = strings.Replace(html, "<!--app-data-->", fmt.Sprintf(DataFormat, data), 1)

			return
		}

		err = errors.New("unknown render mode")

		return
	}
}
