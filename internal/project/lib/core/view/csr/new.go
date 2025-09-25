package csr

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
	_view "github.com/razshare/frizzante/internal/project/lib/core/view"
)

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

	var dist = filepath.Join(app, "dist")
	var index = filepath.Join(dist, "client", "index.html")
	var indexFixed = strings.ReplaceAll(index, "\\", "/")

	return func(view _view.View) (string, error) {
		var indexData []byte
		var err error

		if !disk && embeds.IsFile(efs, indexFixed) {
			indexData, err = efs.ReadFile(indexFixed)
		} else {
			indexData, err = os.ReadFile(index)
		}

		if err != nil {
			return "", err
		}

		indexString := string(indexData)

		var propsData []byte
		if propsData, err = json.Marshal(_view.Wrap(view)); err != nil {
			return "", err
		}

		indexString = strings.Replace(indexString, "<!--app-head-->", fmt.Sprintf(HeadFormat, view.Title), 1)
		indexString = strings.Replace(indexString, "<!--app-body-->", fmt.Sprintf(BodyFormat, ""), 1)
		indexString = strings.Replace(indexString, "<!--app-data-->", fmt.Sprintf(DataFormat, propsData), 1)

		return indexString, nil
	}
}
