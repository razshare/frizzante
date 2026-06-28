//go:build !use_disk

package csr

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

func New(options Options) renders.Render {
	efs := options.Efs
	var index = filepath.Join("app", "dist", "client", "index.html")
	index = strings.ReplaceAll(index, "\\", "/")
	return func(options renders.Options) (document string, err error) {
		var indexData []byte
		if indexData, err = efs.ReadFile(index); err != nil {
			return
		}
		document = string(indexData)
		var data []byte
		if data, err = json.Marshal(options.Data); err != nil {
			return "", err
		}
		document = strings.Replace(document, "<!--app-head-->", fmt.Sprintf(renders.HeadFormat, options.View.Title), 1)
		document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(renders.BodyFormat, ""), 1)
		document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(renders.DataFormat, data), 1)
		return document, nil
	}
}
