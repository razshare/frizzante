package user

import (
	"os"
	"path/filepath"
)

func FrizzanteHome() (home string, err error) {
	home = os.Getenv("FRIZZANTE_HOME")

	if home == "" {
		var user string
		if user, err = os.UserHomeDir(); err != nil {
			return "", err
		}
		home = filepath.Join(user, ".frizzante")
	}

	return
}
