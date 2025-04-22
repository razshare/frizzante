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

var components = map[string]string{}

type Document struct {
	Render     Render
	Data       map[string]any
	view       string
	efs        embed.FS
	parameters map[string]string
}

var noScriptPattern = regexp.MustCompile(`<script.*>.*</script>`)

type DocumentProps struct {
	View       string            `json:"view"`
	Data       map[string]any    `json:"data"`
	Views      map[string]string `json:"views"`
	Parameters map[string]string `json:"parameters"`
}

// DocumentCreate creates a document with a view.

func DocumentCreate(view string) *Document {
	return &Document{
		view: view,
	}
}

// DocumentCompile compiles a document.
func DocumentCompile(self *Document) (string, error) {
	fileNameIndex := filepath.Join(".dist", "client", ".frizzante", "vite-project", "index.html")

	var indexBytes []byte

	if "1" == os.Getenv("DEV") {
		indexBytesLocal, readError := os.ReadFile(fileNameIndex)
		if readError != nil {
			return "", readError
		}
		indexBytes = indexBytesLocal
	} else {
		indexBytesLocal, readError := self.efs.ReadFile(fileNameIndex)
		if readError != nil {
			return "", readError
		}
		indexBytes = indexBytesLocal
	}

	routerPropsBytes, jsonError := json.Marshal(DocumentProps{
		Views:      components,
		View:       self.view,
		Data:       self.Data,
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

	if RenderFull == self.Render {
		head, body, renderError := render(self.efs, routerPropsString)
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

	if RenderClient == self.Render {
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

	if RenderServer == self.Render {
		head, body, renderError := render(self.efs, routerPropsString)
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

	if RenderHeadless == self.Render {
		_, body, renderError := render(self.efs, routerPropsString)

		if renderError != nil {
			return "", renderError
		}

		return body, nil

	}

	return "", nil
}
