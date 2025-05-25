package frizzante

import (
	"embed"
	"encoding/json"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"os"
	"path/filepath"
	"regexp"
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
	Name       string
	Data       any
	Error      error
	RenderMode RenderMode
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
func (view *View) Render(efs embed.FS) (content string, compileError error) {
	var fileNameIndex string

	if "1" == os.Getenv("DEV") {
		fileNameIndex = filepath.Join(".dist", "client", ".frizzante", "vite-project", "index.html")
	} else {
		fileNameIndex = ".dist/client/.frizzante/vite-project/index.html"
	}

	var indexBytes []byte

	if "1" == os.Getenv("DEV") {
		indexBytesLocal, readError := os.ReadFile(fileNameIndex)
		if readError != nil {
			return "", readError
		}
		indexBytes = indexBytesLocal
	} else {
		indexBytesLocal, readError := efs.ReadFile(fileNameIndex)
		if readError != nil {
			return "", readError
		}
		indexBytes = indexBytesLocal
	}

	err := ""
	if nil != view.Error {
		err = view.Error.Error()
	}

	routerPropsBytes, jsonError := json.Marshal(&ServerProperties{
		View:       view.Name,
		RenderMode: view.RenderMode,
		Data:       view.Data,
		Error:      err,
	})

	if jsonError != nil {
		return "", jsonError
	}

	routerPropsString := string(routerPropsBytes)

	targetId, targetIdError := uuid.NewV4()
	if targetIdError != nil {
		return "", targetIdError
	}

	if RenderModeFull == view.RenderMode {
		head, body, renderError := JavaScriptRender(efs, routerPropsString)
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

	if RenderModeClient == view.RenderMode {
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

	if RenderModeServer == view.RenderMode {
		head, body, renderError := JavaScriptRender(efs, routerPropsString)
		if renderError != nil {
			return "", renderError
		}
		return strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						noScript.ReplaceAllString(string(indexBytes), ""),
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

	if RenderModeHeadless == view.RenderMode {
		_, body, renderError := JavaScriptRender(efs, routerPropsString)

		if renderError != nil {
			return "", renderError
		}

		return body, nil

	}

	return "", nil
}
