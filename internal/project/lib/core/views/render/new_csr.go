//go:build !dev && no_js_runtime

package render

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

func New(conf Config) Render {
	var efs = conf.Efs
	var app = conf.App

	if app == "" {
		app = "app"
	}

	var index = filepath.Join(app, "dist", "client", "index.html")

	index = strings.ReplaceAll(index, "\\", "/")

	return func(view views.View) (document string, err error) {
		var indexData []byte
		if indexData, err = efs.ReadFile(index); err != nil {
			return
		}

		document = string(indexData)

		var data []byte
		if data, err = json.Marshal(views.NewData(view)); err != nil {
			return "", err
		}

		document = strings.Replace(document, "<!--app-head-->", fmt.Sprintf(HeadFormat, view.Title), 1)
		document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(BodyFormat, ""), 1)
		document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(DataFormat, data), 1)

		return document, nil
	}
}
