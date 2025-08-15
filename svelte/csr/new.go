package csr

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/view"
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

//go:embed props.format
var PropsFormat string

func New(c Config) view.Render {
	var efs = c.Efs
	var app = c.App
	var disk = c.Disk

	if app == "" {
		app = "app"
	}

	var id = "svelte-app"
	var dist = filepath.Join(app, "dist")
	var docn = filepath.Join(dist, "client", "index.html")
	var docnfix = strings.ReplaceAll(docn, "\\", "/")

	return func(v view.View) (string, error) {
		var d []byte
		var err error

		if !disk && embeds.IsFile(efs, docnfix) {
			d, err = efs.ReadFile(docnfix)
		} else {
			d, err = os.ReadFile(docn)
		}

		if err != nil {
			return "", err
		}

		doc := string(d)

		props, merr := json.Marshal(view.Data(v))

		if merr != nil {
			return "", merr
		}

		doc = strings.Replace(doc, "<!--app-target-->", fmt.Sprintf(TargetFormat, id), 1)
		doc = strings.Replace(doc, "<!--app-head-->", fmt.Sprintf(HeadFormat, v.Title), 1)
		doc = strings.Replace(doc, "<!--app-body-->", fmt.Sprintf(BodyFormat, id, ""), 1)
		doc = strings.Replace(doc, "<!--app-props-->", fmt.Sprintf(PropsFormat, props), 1)

		return doc, nil

	}
}
