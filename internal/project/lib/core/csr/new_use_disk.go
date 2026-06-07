//go:build use_disk

package csr

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

func New() renders.Render {
	var index = filepath.Join("app", "dist", "client", "index.html")
	index = strings.ReplaceAll(index, "/", string(filepath.Separator))
	index = strings.ReplaceAll(index, "\\", string(filepath.Separator))
	return func(options renders.RenderOptions) (document string, err error) {
		var indexData []byte
		if indexData, err = os.ReadFile(index); err != nil {
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
