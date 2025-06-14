package frz

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/fs"
	"os"
	"regexp"
	"rogchap.com/v8go"
	"strings"
)

type RenderMode int

const (
	RenderModeFull     RenderMode = 0 // Renders on both the server and the client.
	RenderModeServer   RenderMode = 1 // Renders only on the server.
	RenderModeClient   RenderMode = 2 // Renders only on the client.
	RenderModeHeadless RenderMode = 3 // Renders only on the server and omits the base template.
)

type EsbuildMode int

const (
	EsbuildModeEnvironment EsbuildMode = 0 // Enables esbuild bundling if environment variable "DEV" equals "1".
	EsbuildModeEnabled     EsbuildMode = 1 // Enables esbuild bundling.
	EsbuildModeDisabled    EsbuildMode = 2 // Disables esbuild bundling.
)

type View struct {
	Name        string         `json:"name"`
	Data        map[string]any `json:"data"`
	RenderMode  RenderMode     `json:"renderMode"`
	EsbuildMode EsbuildMode    `json:"esbuildMode"`
	root        string
	functions   map[string]v8go.FunctionCallback
	server      string
	index       string
}

// WithRoot sets the root of the view.
//
// The root of the view should contain `node_modules` and `package.json`.
func (view *View) WithRoot(root string) *View {
	view.root = root
	return view
}

// WithServer sets the server script.
//
// This script must be in CommonJs format and it must declare a local "render" function.
//
// The render function takes a "props" map parameter and returns a RenderOutput.
//
// See https://svelte.dev/docs/svelte/svelte-server#render
func (view *View) WithServer(fileName string) *View {
	view.server = fileName
	return view
}

// WithIndex sets the index html document that will wrap the final render output.
func (view *View) WithIndex(fileName string) *View {
	view.index = fileName
	return view
}

// AddFunction adds a global function to the script's context.
func (view *View) AddFunction(name string, function v8go.FunctionCallback) *View {
	if nil == view.functions {
		view.functions = map[string]v8go.FunctionCallback{}
	}
	view.functions[name] = function
	return view
}

// IndexContents gets the contents of the index html document.
func (view *View) IndexContents(efs embed.FS) ([]byte, error) {
	var index []byte
	var indexReadError error
	if fs.FileExists(view.index) {
		index, indexReadError = os.ReadFile(view.index)
	}

	if indexReadError != nil || index == nil {
		fileNameFixed := strings.ReplaceAll(view.index, "\\", "/")
		if fs.ExistsInEmbeddedFileSystem(efs, fileNameFixed) {
			index, indexReadError = efs.ReadFile(fileNameFixed)
			if indexReadError != nil {
				return nil, indexReadError
			}
		} else {
			return nil, errors.New("view index is missing from the host file system and the embedded file system")
		}
	}

	return index, nil
}

// ServerContents gets the contents of the server script.
func (view *View) ServerContents(efs embed.FS) ([]byte, error) {
	var server []byte
	var serverReadError error
	if fs.FileExists(view.server) {
		server, serverReadError = os.ReadFile(view.server)
	}

	if serverReadError != nil || server == nil {
		fileNameFixed := strings.ReplaceAll(view.server, "\\", "/")
		if fs.ExistsInEmbeddedFileSystem(efs, fileNameFixed) {
			server, serverReadError = efs.ReadFile(fileNameFixed)
			if serverReadError != nil {
				return nil, serverReadError
			}
		} else {
			return nil, errors.New("view server is missing from the host file system and the embedded file system")
		}
	}
	return server, nil
}

var noScript = regexp.MustCompile(`<script.*>.*</script>`)

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
// If for some reason your view server is not in cjs format, use WithEsbuildMode to enable esbuild and convert the view server script to cjs on the fly.
func (view *View) Render(efs embed.FS) (html string, renderError error) {
	// CSR.
	targetId, targetIdError := uuid.NewV4()
	if targetIdError != nil {
		return "", targetIdError
	}

	propsBytes, marshalError := json.Marshal(view)
	if marshalError != nil {
		return "", marshalError
	}

	props := string(propsBytes)

	if RenderModeClient == view.RenderMode {
		index, indexError := view.IndexContents(efs)
		if indexError != nil {
			return "", indexError
		}
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						string(index),
						"<!--app-target-->",
						fmt.Sprintf("<script type=\"application/javascript\">function target(){return document.getElementById(\"%s\")}</script>", targetId),
						1,
					),
					"<!--app-body-->",
					fmt.Sprintf("<div id=\"%s\"></div>", targetId),
					1,
				),
				"<!--app-head-->",
				"",
				1,
			),
			"<!--app-data-->",
			fmt.Sprintf(
				"<script type=\"application/javascript\">function props(){return %s}</script>",
				props,
			),
			1,
		), nil
	}

	readBytes, serverReadError := view.ServerContents(efs)
	if serverReadError != nil {
		return "", serverReadError
	}

	var serverError error
	var server []byte
	convertToCjs :=
		view.EsbuildMode == EsbuildModeEnabled ||
			(view.EsbuildMode == EsbuildModeEnvironment && os.Getenv("DEV") == "1")

	if convertToCjs {
		server, serverError = JavaScriptBundle(view.root, api.FormatCommonJS, readBytes)
		if serverError != nil {
			return "", serverError
		}
	} else {
		server = readBytes
	}

	serverIif := fmt.Sprintf(
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
		string(server),
	)

	bundle := []byte(serverIif)

	var head string
	var body string
	var err string

	globals := map[string]v8go.FunctionCallback{
		"error": func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			args := info.Args()
			if len(args) > 0 {
				err = err + args[0].String() + "\n"
			}
			return nil
		},
		"props": func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			value, valueError := v8go.NewValue(info.Context().Isolate(), props)
			if valueError != nil {
				return nil
			}
			return value
		},
		"head": func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			args := info.Args()
			if len(args) > 0 {
				head = args[0].String()
			}
			return nil
		},
		"body": func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			args := info.Args()
			if len(args) > 0 {
				body = args[0].String()
			}
			return nil
		},
	}

	if nil != view.functions {
		for name, function := range view.functions {
			globals[name] = function
		}
	}

	_, destroy, javaScriptError := JavaScriptRun(view.server, bundle, globals)
	defer destroy()
	if javaScriptError != nil {
		return "", javaScriptError
	}

	if "" != err {
		return "", errors.New(err)
	}

	if RenderModeHeadless == view.RenderMode {
		return body, nil
	}

	if RenderModeServer == view.RenderMode {
		index, indexError := view.IndexContents(efs)
		if indexError != nil {
			return "", indexError
		}
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						noScript.ReplaceAllString(string(index), ""),
						"<!--app-target-->",
						"",
						1,
					),
					"<!--app-body-->",
					fmt.Sprintf("<div id=\"%s\">%s</div>", targetId, body),
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

	if RenderModeFull == view.RenderMode {
		index, indexError := view.IndexContents(efs)
		if indexError != nil {
			return "", indexError
		}
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						string(index),
						"<!--app-target-->",
						fmt.Sprintf("<script type=\"application/javascript\">function target(){return document.getElementById(\"%s\")}</script>", targetId),
						1,
					),
					"<!--app-body-->",
					fmt.Sprintf("<div id=\"%s\">%s</div>", targetId, body),
					1,
				),
				"<!--app-head-->",
				head,
				1,
			),
			"<!--app-data-->",
			fmt.Sprintf(
				"<script type=\"application/javascript\">function props(){return %s}</script>",
				props,
			),
			1,
		), nil
	}

	return
}
