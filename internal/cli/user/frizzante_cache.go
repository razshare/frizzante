package user

import (
	"path/filepath"
)

func FrizzanteCache() (cache string, err error) {
	var home string
	if home, err = FrizzanteHome(); err != nil {
		return
	}
	cache = filepath.Join(home, "cache")
	return
}
