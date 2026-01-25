package caches

import (
	"fmt"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Install(options InstallOptions) (err error) {
	ext := filepath.Ext(options.FromFileName)
	spin := spinners.New(fmt.Sprintf("installing %s", options.ToDirectoryName))
	go spinners.Start(spin)
	if ext == ".zip" {
		if err = files.UnzipFile(options.FromFileName, options.ToDirectoryName); err != nil {
			spinners.Stop(spin)
			return
		}
	} else {
		local := filepath.Join(options.ToDirectoryName, filepath.Base(options.ToDirectoryName)+ext)
		if err = files.CopyFile(options.FromFileName, local); err != nil {
			spinners.Stop(spin)
			return
		}
	}
	spinners.Stop(spin)
	messages.Successf("%s installed", options.ToDirectoryName)
	return
}
