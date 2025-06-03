package frizzante

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	uuid "github.com/nu7hatch/gouuid"
	"os"
	"path/filepath"
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

type View struct {
	Name       string     `json:"name"`
	Data       any        `json:"data"`
	Error      string     `json:"error"`
	RenderMode RenderMode `json:"renderMode"`
}

var noScript = regexp.MustCompile(`<script.*>.*</script>`)
var bundle []byte

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
func (view *View) Render(efs embed.FS) (html string, jsError error) {
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

	var appBytes []byte

	if "1" == os.Getenv("DEV") {
		var readError error
		appBytes, readError = os.ReadFile(filepath.Join(".dist", "client", ".frz", "router", "index.html"))
		if readError != nil {
			return "", readError
		}
	} else {
		var readError error
		appBytes, readError = efs.ReadFile(".dist/client/.frz/router/index.html")
		if readError != nil {
			return "", readError
		}
	}

	if RenderModeClient == view.RenderMode {
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						string(appBytes),
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

	// SSR.
	if bundle == nil {
		var serverEsmBytes []byte
		if "1" == os.Getenv("DEV") {
			renderFileName := filepath.Join(".dist", "server", "server.js")
			renderEsmBytesLocal, readError := os.ReadFile(renderFileName)
			if readError != nil {
				return "", readError
			}
			serverEsmBytes = renderEsmBytesLocal
		} else {
			renderFileName := ".dist/server/server.js"
			renderEsmBytesLocal, readError := efs.ReadFile(renderFileName)
			if readError != nil {
				return "", readError
			}
			serverEsmBytes = renderEsmBytesLocal
		}

		serverCjsBytes, javaScriptBundleError := JavaScriptBundle(".", api.FormatCommonJS, serverEsmBytes)
		if javaScriptBundleError != nil {
			return "", javaScriptBundleError
		}

		serverIif := fmt.Sprintf(
			`
		const module={exports:{}}; const render = (function(){
			%s
			return render;
		})()
		render(JSON.parse(props())).then(function done(rendered){
			head(rendered.head??'');
			body(rendered.body??'');
		});
		`,
			serverCjsBytes,
		)

		bundleLocal, bundleError := JavaScriptBundle(".", api.FormatCommonJS, []byte(serverIif))
		if bundleError != nil {
			return "", bundleError
		}
		bundle = bundleLocal
	}

	var head string
	var body string

	globals := map[string]v8go.FunctionCallback{
		"props": func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			value, valueError := v8go.NewValue(info.Context().Isolate(), props)
			if nil != valueError {
				return nil
			}
			return value
		},
		"inspect": func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			args := info.Args()
			if len(args) > 0 {
				message := args[0].String()
				println(message)
			}
			return nil
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

	_, destroy, javaScriptError := JavaScriptRun("server.js", bundle, globals)
	defer destroy()
	if javaScriptError != nil {
		return "", javaScriptError
	}

	if RenderModeHeadless == view.RenderMode {
		return body, nil
	}

	if RenderModeServer == view.RenderMode {
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						noScript.ReplaceAllString(string(appBytes), ""),
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
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						string(appBytes),
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
