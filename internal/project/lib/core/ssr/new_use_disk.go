//go:build use_disk

package ssr

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/internal/project/lib/core/esbuild"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/javascript"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

func New(options Options) renders.Render {
	errorLog := options.ErrorLog
	infoLog := options.InfoLog
	if errorLog == nil {
		errorLog = log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime)
	}
	if infoLog == nil {
		infoLog = log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime)
	}
	var server = filepath.Join("app", "dist", "server", "app.server.js")
	var index = filepath.Join("app", "dist", "client", "index.html")
	server = strings.ReplaceAll(server, "/", string(filepath.Separator))
	server = strings.ReplaceAll(server, "\\", string(filepath.Separator))
	index = strings.ReplaceAll(index, "/", string(filepath.Separator))
	index = strings.ReplaceAll(index, "\\", string(filepath.Separator))
	var goRender = func(options renders.Options) (jsRender javascript.Render, err error) {
		if !files.IsFile(server) {
			err = fmt.Errorf("file %s not found", server)
			return
		}
		emptyTsFileName := fmt.Sprintf(".%s%s", string(os.PathSeparator), filepath.Join(".gen", "empty.ts"))
		jsRender, err = javascript.NewRender(javascript.NewRenderOptions{
			Server:   server,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
			FindSource: func() (source string, err error) {
				var data []byte
				if data, err = os.ReadFile(server); err != nil {
					return
				}
				if !files.IsFile(emptyTsFileName) {
					if err = os.WriteFile(emptyTsFileName, []byte("export default {}"), os.ModePerm); err != nil {
						return
					}
				}
				// currently svelte imports `node:crypto` and `node:async_hooks`, which will break our runtime,
				// so we need to strip them off from the bundle.
				// see issues #17762 and #17771:
				// https://github.com/sveltejs/svelte/issues/17762
				// https://github.com/sveltejs/svelte/issues/17771
				source = strings.Replace(string(data), `await obfuscated_import(`, "await import(", 1)
				if source, err = esbuild.Bundle("app", api.FormatCommonJS, source, map[string]string{
					"node:crypto":      strings.ReplaceAll(emptyTsFileName, string(os.PathSeparator), "/"),
					"node:async_hooks": strings.ReplaceAll(emptyTsFileName, string(os.PathSeparator), "/"),
				}); err != nil {
					return
				}
				return
			},
		})
		return
	}
	return func(options renders.Options) (document string, err error) {
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
