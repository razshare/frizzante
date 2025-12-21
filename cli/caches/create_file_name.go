package caches

import (
	"path/filepath"

	"github.com/razshare/frizzante/cli/paths"
	"github.com/razshare/frizzante/internal/additions/lib/security"
)

func CreateFileName(options CreateFileNameOptions) (fileName string, err error) {
	var cache string
	if cache, err = paths.Cache(); err != nil {
		return
	}

	hash := security.Sha1(options.Value)

	ext := filepath.Ext(options.Value)

	if ext != ".exe" && ext != ".zip" {
		ext = ""
	}

	fileName = filepath.Join(cache, hash+ext)
	return
}
