package ssr

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
	"github.com/razshare/frizzante/internal/project/lib/core/js"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	_view "github.com/razshare/frizzante/internal/project/lib/core/view"
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
	if conf.ErrorLog == nil {
		conf.ErrorLog = log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime)
	}
	if conf.InfoLog == nil {
		conf.InfoLog = log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime)
	}

	if limit <= 0 {
		limit = 1
	}

	if app == "" {
		app = "app"
	}

	var mut sync.Mutex
	var id = "app"
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

		var builder strings.Builder
		runtime = goja.New()
		console := runtime.NewObject()
		createLogger := func(level LogLevel) func(call goja.FunctionCall) goja.Value {
			var logger *log.Logger

			switch level {
			case LogLevelDanger:
				logger = conf.ErrorLog
			default:
				logger = conf.InfoLog
			}

			return func(call goja.FunctionCall) goja.Value {
				builder.Reset()
				i := 0
				for _, argument := range call.Arguments {
					if i > 0 {
						builder.WriteString(" ")
					}
					switch argument.(type) {
					case *goja.Object:
						object := argument.ToObject(runtime)
						data, err = object.MarshalJSON()
						if err != nil {
							conf.ErrorLog.Println(err, stack.Trace())
							return goja.Undefined()
						}
						builder.WriteString(string(data))
					default:
						value := argument.String()
						if value == "https://svelte.dev/e/experimental_async_ssr" {
							// Skipping experimental async ssr warnings.
							return goja.Undefined()
						}
						builder.WriteString(value)
					}
					i++
				}
				logger.Println(builder.String())
				return goja.Undefined()
			}
		}

		if err = console.Set("log", createLogger(LogLevelBase)); err != nil {
			return
		}

		if err = console.Set("info", createLogger(LogLevelBase)); err != nil {
			return
		}

		if err = console.Set("warn", createLogger(LogLevelWarning)); err != nil {
			return
		}

		if err = console.Set("error", createLogger(LogLevelDanger)); err != nil {
			return
		}

		if err = runtime.Set("console", console); err != nil {
			return
		}

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
			return
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
			if promise, err = render(goja.Undefined(), runtime.ToValue(_view.Wrap(view))); err != nil {
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
				if data, err = json.Marshal(_view.Wrap(view)); err != nil {
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
			if data, err = json.Marshal(_view.Wrap(view)); err != nil {
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
