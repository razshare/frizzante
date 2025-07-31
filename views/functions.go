package views

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/js"
	"os"
	"rogchap.com/v8go"
	"strings"
)

// IndexHtmlData gets the contents of the index html document.
func IndexHtmlData(self *View, efs embed.FS) ([]byte, error) {
	if files.IsFile(self.IndexHtml) {
		return os.ReadFile(self.IndexHtml)
	}

	var data []byte
	fileNameFixed := strings.ReplaceAll(self.IndexHtml, "\\", "/")
	if embeds.IsFile(efs, fileNameFixed) {
		var readError error
		data, readError = efs.ReadFile(fileNameFixed)
		if readError != nil {
			return nil, readError
		}
	} else {
		return nil, errors.New("view index is missing from the host file system and the embedded file system")
	}
	return data, nil
}

// ServerJsData gets the contents of the server script.
func ServerJsData(self *View, efs embed.FS) ([]byte, error) {
	if files.IsFile(self.ServerJs) {
		return os.ReadFile(self.ServerJs)
	}

	var data []byte
	n := strings.ReplaceAll(self.ServerJs, "\\", "/")
	if embeds.IsFile(efs, n) {
		var readError error
		data, readError = efs.ReadFile(n)
		if readError != nil {
			return nil, readError
		}
	} else {
		return nil, errors.New("view server is missing from the host file system and the embedded file system")
	}
	return data, nil
}

// Render renders the view.
//
// If the View is using RenderModeServer, then ViewRender returns an HTML document.
// The head of the document will contain *only* content declared with the <svelte:head> tag.
// The body of the document will contain the fully rendered content of the view as HTML.
//
// If the View is using RenderModeClient, then ViewRender returns an HTML document.
// The document itself doesn't include any of the view content, instead, custom <script> tags are injected into the head of
// the document in order to asynchronously load a client JavaScript bundle that renders the view inside the client's browser,
// thus ultimately loading the content into the document.
//
// If the View is using RenderModeFull, then ViewRender returns an HTML document.
// The head of the document will contain any content declared with the <svelte:head> tag.
// The body of the document will contain the fully rendered content of the view as HTML.
// On top of that, just like when using RenderModeClient, custom <script> tags are also injected into the head of
// the document in order to asynchronously load a client JavaScript bundle that, in this case, re-renders the
// view inside the client's browser.
// In short, RenderModeFull is a combination of RenderModeServer and RenderModeClient.
//
// If the View is using RenderModeHeadless, then ViewRender returns only the content of the view, without decorating it with an HTML document.
// The output won't even contain a header, ignoring all <svelte:head> declarations and all css.
//
// When rendering the view on the server, the view server (which you can set with WithServer), is expected to be in common js format (cjs).
//
// If for some reason your view server is not in cjs format, RenderMode will try to convert it to cjs on the fly using esbuild.
// Esbuild will look for a "node_modules" in the view root directory, which you can set by invoking WithRoot.
func Render(self *View, efs embed.FS) (html string, err error) {
	// CSR.
	idObject, idObjectError := uuid.NewV4()
	if idObjectError != nil {
		return "", idObjectError
	}

	viewName := self.Name
	viewData := self.Data
	viewRenderMode := self.RenderMode

	if viewData == nil {
		viewData = map[string]any{}
	}

	jsonData, jsonError := json.Marshal(map[string]any{
		"name":       viewName,
		"data":       viewData,
		"renderMode": viewRenderMode,
	})
	if jsonError != nil {
		return "", jsonError
	}

	properties := string(jsonData)

	if RenderModeClient == self.RenderMode {
		indexHtmlData, indexHtmlDataError := IndexHtmlData(self, efs)
		if indexHtmlDataError != nil {
			return "", indexHtmlDataError
		}
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						string(indexHtmlData),
						"<!--app-target-->",
						fmt.Sprintf("<script type=\"application/javascript\">function target(){return document.getElementById(\"%s\")}</script>", idObject),
						1,
					),
					"<!--app-body-->",
					fmt.Sprintf("<div id=\"%s\"></div>", idObject),
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

	var serverJsData []byte
	var serverJsDataError error

	serverJsData, serverJsDataError = ServerJsData(self, efs)
	if serverJsDataError != nil {
		return "", serverJsDataError
	}

	if files.IsDirectory(self.AppRoot) {
		serverJsData, serverJsDataError = js.JavaScriptBundle(self.AppRoot, api.FormatCommonJS, serverJsData)
		if serverJsDataError != nil {
			return "", serverJsDataError
		}
	}

	iif := fmt.Sprintf(
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
		string(serverJsData),
	)

	iifData := []byte(iif)

	var head string
	var body string
	var jsError string

	globalFunctions := map[string]v8go.FunctionCallback{
		"error": func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			arguments := info.Args()
			if len(arguments) > 0 {
				jsError = jsError + arguments[0].String() + "\n"
			}
			return nil
		},
		"props": func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			value, valueError := v8go.NewValue(info.Context().Isolate(), properties)
			if valueError != nil {
				return nil
			}
			return value
		},
		"head": func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			arguments := info.Args()
			if len(arguments) > 0 {
				head = arguments[0].String()
			}
			return nil
		},
		"body": func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			arguments := info.Args()
			if len(arguments) > 0 {
				body = arguments[0].String()
			}
			return nil
		},
	}

	if self.Functions != nil {
		for functionName, functionCallback := range self.Functions {
			globalFunctions[functionName] = functionCallback
		}
	}

	_, destroy, javaScriptError := js.JavaScriptRun(self.ServerJs, iifData, globalFunctions)
	if javaScriptError != nil {
		return "", javaScriptError
	}
	defer destroy()

	if "" != jsError {
		return "", errors.New(jsError)
	}

	if RenderModeHeadless == self.RenderMode {
		return body, nil
	}

	if RenderModeServer == self.RenderMode {
		indexHtmlData, indexHtmlDataError := IndexHtmlData(self, efs)
		if indexHtmlDataError != nil {
			return "", indexHtmlDataError
		}
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						globals.NoScript.ReplaceAllString(string(indexHtmlData), ""),
						"<!--app-target-->",
						"",
						1,
					),
					"<!--app-body-->",
					fmt.Sprintf("<div id=\"%s\">%s</div>", indexHtmlData, body),
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

	if RenderModeFull == self.RenderMode {
		indexHtmlData, indexHtmlDataError := IndexHtmlData(self, efs)
		if indexHtmlDataError != nil {
			return "", indexHtmlDataError
		}
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						string(indexHtmlData),
						"<!--app-target-->",
						fmt.Sprintf("<script type=\"application/javascript\">function target(){return document.getElementById(\"%s\")}</script>", idObject),
						1,
					),
					"<!--app-body-->",
					fmt.Sprintf("<div id=\"%s\">%s</div>", idObject, body),
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

	return
}
