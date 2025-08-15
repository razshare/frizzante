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

//go:embed props.format
var PropsFormat string

var NoScript = regexp.MustCompile(`<script.*>.*</script>`)

func New(c Config) view.Render {
	var efs = c.Efs
	var app = c.App
	var disk = c.Disk
	var limit = c.Limit

	if limit == 0 {
		limit = 1
	}

	if app == "" {
		app = "app"
	}

	var mut sync.Mutex
	var id = "svelte-app"
	var dist = filepath.Join(app, "dist")
	var scriptn = filepath.Join(dist, "server.js")
	var scriptnfix = strings.ReplaceAll(scriptn, "\\", "/")
	var docn = filepath.Join(dist, "client", "index.html")
	var docnfix = strings.ReplaceAll(docn, "\\", "/")
	var renders = make(chan goja.Callable, 1)
	var runtimes = make(chan *goja.Runtime, 1)
	var compile = func() (render goja.Callable, r *goja.Runtime, err error) {
		var d []byte

		if !disk && embeds.IsFile(efs, scriptnfix) {
			d, err = efs.ReadFile(scriptnfix)
		} else {
			d, err = os.ReadFile(scriptn)
		}

		if err != nil {
			return
		}

		r = goja.New()

		var src string
		src, err = js.Bundle(app, api.FormatCommonJS, string(d))
		if err != nil {
			return
		}

		var prog *goja.Program
		prog, err = goja.Compile(scriptn, fmt.Sprintf(RenderFormat, src), false)
		if err != nil {
			return
		}

		var v goja.Value

		v, err = r.RunProgram(prog)
		if err != nil {
			return
		}

		var isf bool
		render, isf = goja.AssertFunction(v)

		if !isf {
			err = errors.New("render is not a function")
			return
		}
		return
	}

	return func(v view.View) (string, error) {
		var d []byte
		var err error

		if !disk && embeds.IsFile(efs, docnfix) {
			d, err = efs.ReadFile(docnfix)
		} else {
			d, err = os.ReadFile(docn)
		}

		if err != nil {
			return "", err
		}

		doc := string(d)

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
				doc = NoScript.ReplaceAllString(doc, "")
			}

			if v.Render == view.RenderServer {
				doc = strings.Replace(doc, "<!--app-target-->", "", 1)
				doc = strings.Replace(doc, "<!--app-data-->", "", 1)
			} else {
				props, merr := json.Marshal(view.Data(v))
				if merr != nil {
					return "", merr
				}
				doc = strings.Replace(doc, "<!--app-target-->", fmt.Sprintf(TargetFormat, id), 1)
				doc = strings.Replace(doc, "<!--app-data-->", fmt.Sprintf(PropsFormat, props), 1)
			}

			doc = strings.Replace(doc, "<!--app-head-->", head, 1)
			doc = strings.Replace(doc, "<!--app-body-->", fmt.Sprintf(BodyFormat, id, body), 1)

			return doc, nil
		}

		if v.Render == view.RenderClient {
			props, merr := json.Marshal(view.Data(v))

			if merr != nil {
				return "", merr
			}

			doc = strings.Replace(doc, "<!--app-target-->", fmt.Sprintf(TargetFormat, id), 1)
			doc = strings.Replace(doc, "<!--app-body-->", fmt.Sprintf(BodyFormat, id, ""), 1)
			doc = strings.Replace(doc, "<!--app-head-->", fmt.Sprintf(HeadFormat, v.Title), 1)
			doc = strings.Replace(doc, "<!--app-data-->", fmt.Sprintf(PropsFormat, props), 1)

			return doc, nil
		}

		return "", errors.New("unknown render mode")
	}
}
