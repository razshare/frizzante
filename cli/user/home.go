package user

import (
	"os"
	"path/filepath"
)

func FrizzanteHome() (string, error) {
	home := os.Getenv("FRIZZANTE_HOME")
	if home == "" {
		var user string
		user, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		home = filepath.Join(user, ".frizzante")
	}

	return home, nil
}
