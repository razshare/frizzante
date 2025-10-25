package paths

import (
	"path/filepath"
)

func Cache() (cache string, err error) {
	var home string
	if home, err = Home(); err != nil {
		return
	}
	cache = filepath.Join(home, "cache")
	return
}
