package csr

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_view "main/lib/core/view"

	"github.com/razshare/frizzante/embeds"
)

//go:embed target.format
var TargetFormat string

//go:embed head.format
var HeadFormat string

//go:embed body.format
var BodyFormat string

//go:embed data.format
var DataFormat string

func New(conf Config) func(view _view.View) (html string, err error) {
	var efs = conf.Efs
	var app = conf.App
	var disk = conf.Disk

	if app == "" {
		app = "app"
	}

	var id = "svelte-app"
	var dist = filepath.Join(app, "dist")
	var index = filepath.Join(dist, "client", "index.html")
	var indexFixed = strings.ReplaceAll(index, "\\", "/")

	return func(view _view.View) (string, error) {
		var data []byte
		var err error

		if !disk && embeds.IsFile(efs, indexFixed) {
			data, err = efs.ReadFile(indexFixed)
		} else {
			data, err = os.ReadFile(index)
		}

		if err != nil {
			return "", err
		}

		doc := string(data)

		var props []byte
		if props, err = json.Marshal(_view.Data(view)); err != nil {
			return "", err
		}

		doc = strings.Replace(doc, "<!--app-target-->", fmt.Sprintf(TargetFormat, id), 1)
		doc = strings.Replace(doc, "<!--app-head-->", fmt.Sprintf(HeadFormat, view.Title), 1)
		doc = strings.Replace(doc, "<!--app-body-->", fmt.Sprintf(BodyFormat, id, ""), 1)
		doc = strings.Replace(doc, "<!--app-props-->", fmt.Sprintf(DataFormat, props), 1)

		return doc, nil
	}
}
