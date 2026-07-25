package caches

import (
	"fmt"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/tui/spinners"
)

func DownloadFile(options DownloadFileOptions) (fileName string, err error) {
	if fileName, err = CreateFileName(CreateFileNameOptions{Value: options.Url}); err != nil {
		return
	}
	if !files.IsFile(fileName) {
		spin := spinners.New(fmt.Sprintf("downloading %s", options.Url))
		go spinners.Start(spin)
		if err = files.DownloadFile(options.Url, fileName); err != nil {
			spinners.Stop(spin)
			return
		}
		spinners.Stop(spin)
	}
	return
}
