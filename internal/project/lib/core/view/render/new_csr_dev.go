//go:build dev && no_js_runtime && !experimental_qjs_runtime

package render

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	view_ "github.com/razshare/frizzante/internal/project/lib/core/view"
)

func New(conf Config) Render {
	var app = conf.App

	if app == "" {
		app = "app"
	}

	var index = filepath.Join(app, "dist", "client", "index.html")

	index = strings.ReplaceAll(index, "/", string(filepath.Separator))
	index = strings.ReplaceAll(index, "\\", string(filepath.Separator))

	return func(view view_.View) (document string, err error) {
		if !files.IsFile(index) {
			err = fmt.Errorf("file %s not found", index)
			return
		}

		var indexData []byte
		if indexData, err = os.ReadFile(index); err != nil {
			return
		}

		document = string(indexData)

		var data []byte
		if data, err = json.Marshal(view_.NewData(view)); err != nil {
			return "", err
		}

		document = strings.Replace(document, "<!--app-head-->", fmt.Sprintf(HeadFormat, view.Title), 1)
		document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(BodyFormat, ""), 1)
		document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(DataFormat, data), 1)

		return document, nil
	}
}
