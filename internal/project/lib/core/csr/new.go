//go:build !use_disk

package csr

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

func New(options Options) renders.Render {
	efs := options.Efs
	var index = filepath.Join("app", "dist", "client", "index.html")
	index = strings.ReplaceAll(index, "\\", "/")
	return func(options renders.Options) (document string, err error) {
		if options.View.RenderMode == views.RenderModeServer {
			err = errors.New("exclusive server rendering is not allowed using a csr function")
			return
		}
		var indexData []byte
		if indexData, err = efs.ReadFile(index); err != nil {
			return
		}
		document = string(indexData)
		var data []byte
		if data, err = json.Marshal(options.Data); err != nil {
			return "", err
		}
		document = strings.Replace(document, "<!--app-head-->", fmt.Sprintf(renders.HeadFormat, ""), 1)
		document = strings.Replace(document, "<!--app-body-->", fmt.Sprintf(renders.BodyFormat, ""), 1)
		document = strings.Replace(document, "<!--app-data-->", fmt.Sprintf(renders.DataFormat, data), 1)
		return document, nil
	}
}
