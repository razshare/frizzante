package views

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/javascript"
	"os"
	"strings"
	"sync"
)

// ReadIndexHtml reads the contents of the index html document and returns it.
func (view *View) ReadIndexHtml(efs embed.FS) (string, error) {
	if files.IsFile(view.IndexHtml) {
		data, readError := os.ReadFile(view.IndexHtml)
		if readError != nil {
			return "", readError
		}
		return string(data), nil
	}

	var data []byte
	fileNameFixed := strings.ReplaceAll(view.IndexHtml, "\\", "/")
	if embeds.IsFile(efs, fileNameFixed) {
		var readError error
		data, readError = efs.ReadFile(fileNameFixed)
		if readError != nil {
			return "", readError
		}
	} else {
		return "", errors.New("view index is missing from the host file system and the embedded file system")
	}
	return string(data), nil
}

// ReadServerJs reads the contents of the server script and returns it.
func (view *View) ReadServerJs(efs embed.FS) (*goja.Program, string, error) {
	fix := func(sourceCode string) string {
		return fmt.Sprintf(
			`
			const module={exports:{}}; const render = (function(){
				%s
				return render;
			})()
			render(JSON.parse(props())).then(function success(r){
				head(r.head??'');
				body(r.body??'');
			}).catch(function failure(e){
				error(e.stack)
			});
			`,
			sourceCode,
		)

	}

	if files.IsFile(view.ServerJs) {
		data, readError := os.ReadFile(view.ServerJs)
		if readError != nil {
			return nil, "", readError
		}

		if files.IsDirectory(view.AppRoot) {
			sourceCode, bundleError := javascript.Bundle(view.AppRoot, api.FormatCommonJS, string(data))
			if bundleError != nil {
				return nil, "", bundleError
			}

			return nil, fix(sourceCode), nil
		}

		if Program == nil {
			Mutex.Lock()
			var compileError error
			Program, compileError = goja.Compile("goja", fix(string(data)), false)
			if compileError != nil {
				Mutex.Unlock()
				return nil, "", compileError
			}
			Mutex.Unlock()
		}

		return Program, "", nil
	}

	fileName := strings.ReplaceAll(view.ServerJs, "\\", "/")

	if embeds.IsFile(efs, fileName) {
		data, readError := efs.ReadFile(fileName)
		if readError != nil {
			return nil, "", readError
		}

		if files.IsDirectory(view.AppRoot) {
			sourceCode, bundleError := javascript.Bundle(view.AppRoot, api.FormatCommonJS, string(data))
			if bundleError != nil {
				return nil, "", bundleError
			}

			return nil, fix(sourceCode), nil
		}

		if Program == nil {
			Mutex.Lock()
			var compileError error
			Program, compileError = goja.Compile("goja", fix(string(data)), false)
			if compileError != nil {
				Mutex.Unlock()
				return nil, "", compileError
			}
			Mutex.Unlock()
		}

		return Program, "", nil
	}

	return nil, "", errors.New("view server is missing from the host file system and the embedded file system")
}

func (view *View) StringifyProperties() (string, error) {
	jsonData, jsonError := json.Marshal(map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if jsonError != nil {
		return "", jsonError
	}

	return string(jsonData), nil
}

var Program *goja.Program
var Mutex sync.Mutex

// ExecuteServerJs executes the server script.
func (view *View) ExecuteServerJs(efs embed.FS) (head string, body string, properties string, err error) {
	var jsError string

	properties, propertiesError := view.StringifyProperties()

	if propertiesError != nil {
		return "", "", "", propertiesError
	}

	program, sourceCode, readError := view.ReadServerJs(efs)
	if readError != nil {
		return "", "", "", readError
	}

	runtime := javascript.New()

	setError := runtime.SetFunctions(map[string]javascript.Function{
		"error": func(call goja.FunctionCall) goja.Value {
			arguments := call.Arguments
			if len(arguments) > 0 {
				jsError = jsError + arguments[0].String() + "\n"
			}
			return nil
		},
		"props": func(call goja.FunctionCall) goja.Value {
			return runtime.ToValue(properties)
		},
		"head": func(call goja.FunctionCall) goja.Value {
			arguments := call.Arguments
			if len(arguments) > 0 {
				head = arguments[0].String()
			}
			return nil
		},
		"body": func(call goja.FunctionCall) goja.Value {
			arguments := call.Arguments
			if len(arguments) > 0 {
				body = arguments[0].String()
			}
			return nil
		},
	})

	if setError != nil {
		return "", "", "", setError
	}

	if program != nil {
		_, runError := runtime.RunProgram(program)
		if runError != nil {
			return "", "", "", runError
		}
	} else {
		_, runError := runtime.RunScript("goja", sourceCode)
		if runError != nil {
			return "", "", "", runError
		}
	}

	if "" != jsError {
		return "", "", "", errors.New(jsError)
	}

	return head, body, properties, nil
}

// RenderClient renders on the client.
func (view *View) RenderClient(efs embed.FS) (string, error) {
	id, idError := uuid.NewV4()
	if idError != nil {
		return "", idError
	}

	properties, propertiesError := view.StringifyProperties()
	if propertiesError != nil {
		return "", propertiesError
	}

	indexHtmlData, indexHtmlDataError := view.ReadIndexHtml(efs)
	if indexHtmlDataError != nil {
		return "", indexHtmlDataError
	}

	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					indexHtmlData,
					"<!--app-target-->",
					fmt.Sprintf("<script type=\"application/javascript\">function target(){return document.getElementById(\"%s\")}</script>", id),
					1,
				),
				"<!--app-body-->",
				fmt.Sprintf("<div id=\"%s\"></div>", id),
				1,
			),
			"<!--app-head-->",
			"",
			1,
		),
		"<!--app-data-->",
		fmt.Sprintf(
			"<script type=\"application/javascript\">function props(){return %s}</script>",
			properties,
		),
		1,
	), nil
}

// RenderServer renders on the server.
func (view *View) RenderServer(efs embed.FS) (string, error) {
	head, body, _, err := view.ExecuteServerJs(efs)
	if err != nil {
		return "", err
	}

	index, readError := view.ReadIndexHtml(efs)
	if readError != nil {
		return "", readError
	}

	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					globals.NoScript.ReplaceAllString(index, ""),
					"<!--app-target-->",
					"",
					1,
				),
				"<!--app-body-->",
				fmt.Sprintf("<div id=\"%s\">%s</div>", index, body),
				1,
			),
			"<!--app-head-->",
			head,
			1,
		),
		"<!--app-data-->",
		"",
		1,
	), nil
}

// RenderHeadless renders only the body of the view on the server.
func (view *View) RenderHeadless(efs embed.FS) (string, error) {
	_, body, _, err := view.ExecuteServerJs(efs)
	if err != nil {
		return "", err
	}
	return body, err
}

// RenderFull renders on the server and on the client.
func (view *View) RenderFull(efs embed.FS) (string, error) {
	id, idError := uuid.NewV4()
	if idError != nil {
		return "", idError
	}

	head, body, properties, err := view.ExecuteServerJs(efs)
	if err != nil {
		return "", err
	}

	index, readError := view.ReadIndexHtml(efs)
	if readError != nil {
		return "", readError
	}

	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					index,
					"<!--app-target-->",
					fmt.Sprintf("<script type=\"application/javascript\">function target(){return document.getElementById(\"%s\")}</script>", id),
					1,
				),
				"<!--app-body-->",
				fmt.Sprintf("<div id=\"%s\">%s</div>", id, body),
				1,
			),
			"<!--app-head-->",
			head,
			1,
		),
		"<!--app-data-->",
		fmt.Sprintf(
			"<script type=\"application/javascript\">function props(){return %s}</script>",
			properties,
		),
		1,
	), nil
}

// Render renders.
func (view *View) Render(efs embed.FS) (string, error) {
	if view.RenderMode == RenderModeServer {
		return view.RenderServer(efs)
	}

	if view.RenderMode == RenderModeFull {
		return view.RenderFull(efs)
	}

	if view.RenderMode == RenderModeClient {
		return view.RenderClient(efs)
	}

	return view.RenderHeadless(efs)
}
