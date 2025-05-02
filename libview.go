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

type Render int64

const (
	RenderServer   Render = 0 // Renders only on the server.
	RenderClient   Render = 1 // Renders only on the client.
	RenderFull     Render = 2 // Renders on both the server and the client.
	RenderHeadless Render = 3 // Renders only on the server and omits the base template.
)

var components = map[string]string{}

type View struct {
	render             Render
	data               map[string]any
	functions          map[string]func(info *v8go.FunctionCallbackInfo) *v8go.Value
	name               string
	parameters         map[string]string
	notifier           *Notifier
	embeddedFileSystem *embed.FS
}

var noScriptPattern = regexp.MustCompile(`<script.*>.*</script>`)

type ViewProps struct {
	View       string            `json:"view"`
	Data       map[string]any    `json:"data"`
	Views      map[string]string `json:"views"`
	Parameters map[string]string `json:"parameters"`
}

// ViewReference references a view in lib/components/views.
func ViewReference(view string) *View {
	return &View{
		name:       view,
		data:       map[string]any{},
		parameters: map[string]string{},
		functions:  map[string]func(info *v8go.FunctionCallbackInfo) *v8go.Value{},
	}
}

// ViewWithRender sets the render mode for the view.
func ViewWithRender(self *View, render Render) {
	self.render = render
}

// ViewWithData sets data for the view.
func ViewWithData(self *View, key string, value any) {
	self.data[key] = value
}

// ViewWithNotifier sets the view notifier.
//
// Note that when using ServerWithPage, the notifier of the view falls
// back to the server notifier.
//
// Unless you have some very specific use case, this means you don't need
// to worry about setting the notifier every time you create a new View,
// just set it once on the server using ServerWithNotifier.
func ViewWithNotifier(self *View, notifier *Notifier) {
	self.notifier = notifier
}

// ViewWithFunction sets a global function for the view.
func ViewWithFunction(
	self *View,
	name string,
	function func(info *v8go.FunctionCallbackInfo) *v8go.Value,
) {
	self.functions[name] = function
}

// ViewWithEmbeddedFileSystem sets the embedded file system for the view.
//
// Note that when using ServerWithPage, the embedded file system of the view falls
// back to the server embedded file system.
//
// Unless you have some very specific use case, this means you don't need
// to worry about setting the embedded file system every time you create a new View,
// just set it once on the server using ServerWithEmbeddedFileSystem.
func ViewWithEmbeddedFileSystem(self *View, embeddedFileSystem embed.FS) {
	self.embeddedFileSystem = &embeddedFileSystem
}

// ViewRender renders a view.
//
// If the View is using RenderServer, then ViewRender returns an HTML document.
// The head of the document will contain *only* content declared with the <svelte:head> tag.
// The body of the document will contain the fully rendered content of the view as HTML.
//
// If the View is using RenderClient, then ViewRender returns an HTML document.
// The document itself doesn't include any of the view content, instead, custom <script> tags are injected into the head of
// the document in order to asynchronously load a client JavaScript bundle that renders the view inside the client's browser,
// thus ultimately loading the content into the document.
//
// If the View is using RenderFull, then ViewRender returns an HTML document.
// The head of the document will contain any content declared with the <svelte:head> tag.
// The body of the document will contain the fully rendered content of the view as HTML.
// On top of that, just like when using RenderClient, custom <script> tags are also injected into the head of
// the document in order to asynchronously load a client JavaScript bundle that, in this case, re-renders the
// view inside the client's browser.
// In short, RenderFull is a combination of RenderServer and RenderClient.
//
// If the View is using RenderHeadless, then ViewRender returns only the content of the view, without decorating it with an HTML document.
// The output won't even contain a header, ignoring all <svelte:head> declarations and all css.
func ViewRender(self *View) (content string, compileError error) {
	fileNameIndex := filepath.Join(".dist", "client", ".frizzante", "vite-project", "index.html")

	var indexBytes []byte

	if "1" == os.Getenv("DEV") {
		indexBytesLocal, readError := os.ReadFile(fileNameIndex)
		if readError != nil {
			return "", readError
		}
		indexBytes = indexBytesLocal
	} else {
		indexBytesLocal, readError := self.embeddedFileSystem.ReadFile(fileNameIndex)
		if readError != nil {
			return "", readError
		}
		indexBytes = indexBytesLocal
	}

	routerPropsBytes, jsonError := json.Marshal(ViewProps{
		Views:      components,
		View:       self.name,
		Data:       self.data,
		Parameters: self.parameters,
	})

	if jsonError != nil {
		return "", jsonError
	}

	routerPropsString := string(routerPropsBytes)

	targetId, targetIdError := uuid.NewV4()
	if targetIdError != nil {
		return "", targetIdError
	}

	if RenderFull == self.render {
		head, body, renderError := ViewExecuteRenderServerJs(self, routerPropsString)
		if renderError != nil {
			return "", renderError
		}
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						string(indexBytes),
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
				routerPropsString,
			),
			1,
		), nil
	}

	if RenderClient == self.render {
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						string(indexBytes),
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
				routerPropsString,
			),
			1,
		), nil
	}

	if RenderServer == self.render {
		head, body, renderError := ViewExecuteRenderServerJs(self, routerPropsString)
		if renderError != nil {
			return "", renderError
		}
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						noScriptPattern.ReplaceAllString(string(indexBytes), ""),
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

	if RenderHeadless == self.render {
		_, body, renderError := ViewExecuteRenderServerJs(self, routerPropsString)

		if renderError != nil {
			return "", renderError
		}

		return body, nil

	}

	return "", nil
}

// ViewExecuteRenderServerJs executes the `.dist/server/render.server.js` file
// and returns the head of the document along with its body.
//
// If the environment variable DEV is set to 1, the file .dist/server/render.server.js is executed directly from the
// local file system, otherwise ViewExecuteRenderServerJs executes the file .dist/server/render.server.js located within the
// view's embedded file system.
func ViewExecuteRenderServerJs(self *View, stringifiedProps string) (head string, body string, jsError error) {
	renderFileName := filepath.Join(".dist", "server", "render.server.js")

	var renderEsmBytes []byte
	if "1" == os.Getenv("DEV") {
		renderEsmBytesLocal, readError := os.ReadFile(renderFileName)
		if readError != nil {
			return "", "", readError
		}
		renderEsmBytes = renderEsmBytesLocal
	} else {
		renderEsmBytesLocal, readError := self.embeddedFileSystem.ReadFile(renderFileName)
		if readError != nil {
			return "", "", readError
		}
		renderEsmBytes = renderEsmBytesLocal
	}

	renderEsm := string(renderEsmBytes)

	renderCjs, javaScriptBundleError := JavaScriptBundle(".", api.FormatCommonJS, renderEsm)
	if javaScriptBundleError != nil {
		return "", "", javaScriptBundleError
	}

	renderIif := fmt.Sprintf("const module={exports:{}}; const render = \n(function(){\n%s\nreturn render;\n})()", renderCjs)

	doneEsm := fmt.Sprintf(
		`
		%s
		render(JSON.parse(stringifiedProps())).then(function done(rendered){
			head(rendered.head??'');
			body(rendered.body??'');
		});
		`,
		renderIif,
	)

	doneCjs, bundleError := JavaScriptBundle(".", api.FormatCommonJS, doneEsm)
	if bundleError != nil {
		return "", "", bundleError
	}

	globals := map[string]v8go.FunctionCallback{}

	if nil != self.functions {
		for name, function := range self.functions {
			globals[name] = function
		}
	}

	globals["stringifiedProps"] = func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		value, valueError := v8go.NewValue(info.Context().Isolate(), stringifiedProps)
		if nil != valueError {
			return nil
		}
		return value
	}

	globals["inspect"] = func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) > 0 {
			message := args[0].String()
			println(message)
		}
		return nil
	}

	globals["head"] = func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) > 0 {
			head = args[0].String()
		}
		return nil
	}

	globals["body"] = func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) > 0 {
			body = args[0].String()
		}
		return nil
	}

	_, destroy, javaScriptError := JavaScriptRun(renderFileName, doneCjs, globals)
	defer destroy()
	if javaScriptError != nil {
		return head, body, javaScriptError
	}

	return head, body, nil
}
