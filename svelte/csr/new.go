package csr

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/razshare/frizzante/embeds"
	v "github.com/razshare/frizzante/view"
	"os"
	"path/filepath"
	"strings"
)

//go:embed target.format
var TargetFormat string

//go:embed head.format
var HeadFormat string

//go:embed body.format
var BodyFormat string

//go:embed data.format
var DataFormat string

func New(conf Config) v.Render {
	var efs = conf.Efs
	var app = conf.App
	var disk = conf.Disk

	if app == "" {
		app = "app"
	}

	var id = "svelte-app"
	var nameDist = filepath.Join(app, "dist")
	var nameDoc = filepath.Join(nameDist, "client", "index.html")
	var nameDocFixed = strings.ReplaceAll(nameDoc, "\\", "/")

	return func(view v.View) (string, error) {
		var data []byte
		var err error

		if !disk && embeds.IsFile(efs, nameDocFixed) {
			data, err = efs.ReadFile(nameDocFixed)
		} else {
			data, err = os.ReadFile(nameDoc)
		}

		if err != nil {
			return "", err
		}

		doc := string(data)

		var props []byte
		if props, err = json.Marshal(v.Data(view)); err != nil {
			return "", err
		}

		doc = strings.Replace(doc, "<!--app-target-->", fmt.Sprintf(TargetFormat, id), 1)
		doc = strings.Replace(doc, "<!--app-head-->", fmt.Sprintf(HeadFormat, view.Title), 1)
		doc = strings.Replace(doc, "<!--app-body-->", fmt.Sprintf(BodyFormat, id, ""), 1)
		doc = strings.Replace(doc, "<!--app-props-->", fmt.Sprintf(DataFormat, props), 1)

		return doc, nil
	}
}
