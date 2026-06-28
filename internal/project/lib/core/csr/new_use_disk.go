//go:build use_disk

package csr

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

func New(_ Options) renders.Render {
	var index = filepath.Join("app", "dist", "client", "index.html")
	var baseDocument string
	baseDocument = strings.Replace(baseDocument, "<!--app-body-->", fmt.Sprintf(renders.BodyFormat, ""), 1)
	baseDocument = strings.Replace(baseDocument, "<!--app-head-->", fmt.Sprintf(renders.HeadFormat, ""), 1)
	index = strings.ReplaceAll(index, "/", string(filepath.Separator))
	index = strings.ReplaceAll(index, "\\", string(filepath.Separator))
	return func(options renders.Options) (document string, err error) {
		if options.View.RenderMode == views.RenderModeServer {
			return
		}
		var data []byte
		if data, err = json.Marshal(options.Data); err != nil {
			return
		}
		var indexData []byte
		if indexData, err = os.ReadFile(index); err != nil {
			return
		}
		document = strings.Replace(string(indexData), "<!--app-head-->", fmt.Sprintf(renders.HeadFormat, ""), 1)
		document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(renders.BodyFormat, ""), 1)
		document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(renders.DataFormat, data), 1)
		return document, nil
	}
}
