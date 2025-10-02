//go:build !dev && no_js_runtime

package render

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
	view_ "github.com/razshare/frizzante/internal/project/lib/core/view"
)

func New(conf Config) func(view view_.View) (html string, err error) {
	var efs = conf.Efs
	var app = conf.App

	if app == "" {
		app = "app"
	}

	var index = filepath.Join(app, "dist", "client", "index.html")

	index = strings.ReplaceAll(index, "\\", "/")

	return func(view view_.View) (string, error) {
		if view.RenderMode != view_.RenderModeClient {
			return "", nil
		}

		var indexData []byte
		var err error

		if embeds.IsFile(efs, index) {
			indexData, err = efs.ReadFile(index)
		}

		if err != nil {
			return "", err
		}

		indexString := string(indexData)

		var propsData []byte
		if propsData, err = json.Marshal(view_.NewData(view)); err != nil {
			return "", err
		}

		indexString = strings.Replace(indexString, "<!--app-head-->", fmt.Sprintf(HeadFormat, view.Title), 1)
		indexString = strings.Replace(indexString, "<!--app-body-->", fmt.Sprintf(BodyFormat, ""), 1)
		indexString = strings.Replace(indexString, "<!--app-data-->", fmt.Sprintf(DataFormat, propsData), 1)

		return indexString, nil
	}
}
