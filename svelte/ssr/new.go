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
	var dst = filepath.Join(app, "dist")
	var src = filepath.Join(dst, "server.js")
	var doc = filepath.Join(dst, "client", "index.html")
	var srcfix = strings.ReplaceAll(src, "\\", "/")
	var docfix = strings.ReplaceAll(doc, "\\", "/")
	var rndrs = make(chan goja.Callable, 1)
	var rntms = make(chan *goja.Runtime, 1)
	var comp = func() (rndr goja.Callable, rntm *goja.Runtime, err error) {
		var d []byte

		if !disk && embeds.IsFile(efs, srcfix) {
			d, err = efs.ReadFile(srcfix)
		} else {
			d, err = os.ReadFile(src)
		}

		if err != nil {
			return
		}

		rntm = goja.New()

		var txt string
		if txt, err = js.Bundle(app, api.FormatCommonJS, string(d)); err != nil {
			return
		}

		var prog *goja.Program
		if prog, err = goja.Compile(src, fmt.Sprintf(RenderFormat, txt), false); err != nil {
			return
		}

		var val goja.Value
		if val, err = rntm.RunProgram(prog); err != nil {
			return
		}

		var isfun bool
		if rndr, isfun = goja.AssertFunction(val); !isfun {
			err = errors.New("render is not a function")
		}

		return
	}

	return func(v view.View) (html string, err error) {
		var d []byte

		if !disk && embeds.IsFile(efs, docfix) {
			d, err = efs.ReadFile(docfix)
		} else {
			d, err = os.ReadFile(doc)
		}

		if err != nil {
			return "", err
		}

		html = string(d)

		if v.Render == view.RenderServer || v.Render == view.RenderFull {
			var rndr goja.Callable
			var rntm *goja.Runtime
			if disk {
				rndr, rntm, err = comp()
			} else if limit >= 0 {
				mut.Lock()
				if limit >= 0 {
					limit--
				}
				mut.Unlock()
				rndr, rntm, err = comp()
				defer func() { go func() { rndrs <- rndr }() }()
				defer func() { go func() { rntms <- rntm }() }()
			} else {
				rndr = <-rndrs
				rntm = <-rntms
				defer func() { go func() { rndrs <- rndr }() }()
				defer func() { go func() { rntms <- rntm }() }()
			}

			if err != nil {
				return "", err
			}

			prms, perr := rndr(goja.Undefined(), rntm.ToValue(view.Data(v)))

			if perr != nil {
				return "", perr
			}

			val := prms.Export().(*goja.Promise).Result().ToObject(rntm)

			headv := val.Get("head")
			bodyv := val.Get("body")

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
				if d, err = json.Marshal(view.Data(v)); err != nil {
					return
				}

				html = strings.Replace(html, "<!--app-target-->", fmt.Sprintf(TargetFormat, id), 1)
				html = strings.Replace(html, "<!--app-data-->", fmt.Sprintf(DataFormat, d), 1)
			}

			html = strings.Replace(html, "<!--app-head-->", head, 1)
			html = strings.Replace(html, "<!--app-body-->", fmt.Sprintf(BodyFormat, id, body), 1)

			return html, nil
		}

		if v.Render == view.RenderClient {
			if d, err = json.Marshal(view.Data(v)); err != nil {
				return
			}

			html = strings.Replace(html, "<!--app-target-->", fmt.Sprintf(TargetFormat, id), 1)
			html = strings.Replace(html, "<!--app-body-->", fmt.Sprintf(BodyFormat, id, ""), 1)
			html = strings.Replace(html, "<!--app-head-->", fmt.Sprintf(HeadFormat, v.Title), 1)
			html = strings.Replace(html, "<!--app-data-->", fmt.Sprintf(DataFormat, d), 1)

			return
		}

		err = errors.New("unknown render mode")

		return
	}
}
