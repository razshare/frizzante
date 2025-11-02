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

var Limit int

func New() Render {
	var index = filepath.Join("app", "dist", "client", "index.html")

	index = strings.ReplaceAll(index, "\\", "/")

	return func(options Options) (document string, err error) {
		var indexData []byte
		if indexData, err = options.Efs.ReadFile(index); err != nil {
			return
		}

		document = string(indexData)

		var data []byte
		if data, err = json.Marshal(views.NewData(options.View)); err != nil {
			return "", err
		}

		document = strings.Replace(document, "<!--app-head-->", fmt.Sprintf(HeadFormat, options.View.Title), 1)
		document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(BodyFormat, ""), 1)
		document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(DataFormat, data), 1)

		return document, nil
	}
}
