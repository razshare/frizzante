package cli

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Go(basepath string) string {
	var goBinary string

	if *FlagGo != "" {
		goBinary = *FlagGo
	} else {
		goBinary = Go(".")
	}

	if strings.HasPrefix(goBinary, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		goBinary = strings.Replace(goBinary, "~", dirname, 1)
		return goBinary + Extension()
	}

	if !strings.Contains(goBinary, string(filepath.Separator)) {
		return goBinary + Extension()
	}

	var pathError error
	goBinary, pathError = filepath.Rel(basepath, goBinary)
	if pathError != nil {
		Fatal(pathError)
	}

	return goBinary + Extension()
}

func Air(basepath string) string {
	var air string

	if *FlagAir != "" {
		air = *FlagAir
	} else {
		air = filepath.Join(".gen", "air", "air")
	}

	if strings.HasPrefix(air, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		air = strings.Replace(air, "~", dirname, 1)
		return air + Extension()
	}

	if !strings.Contains(air, string(filepath.Separator)) {
		return air + Extension()
	}

	var pathError error

	air, pathError = filepath.Rel(basepath, air)
	if pathError != nil {
		Fatal(pathError)
	}

	return air + Extension()
}

func Bun(basepath string) string {
	var bun string

	if *FlagBun != "" {
		bun = *FlagBun
	} else {
		bun = filepath.Join(".gen", "bun", "bun")
	}

	if strings.HasPrefix(bun, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		bun = strings.Replace(bun, "~", dirname, 1)
		return bun + Extension()
	}

	if !strings.Contains(bun, string(filepath.Separator)) {
		return bun + Extension()
	}

	var pathError error
	bun, pathError = filepath.Rel(basepath, bun)
	if pathError != nil {
		Fatal(pathError)
	}

	return bun + Extension()
}

func Sqlc(basepath string) string {
	var sqlc string

	if *FlagSqlc != "" {
		sqlc = *FlagSqlc
	} else {
		sqlc = filepath.Join(".gen", "sqlc", "sqlc")
	}

	if strings.HasPrefix(sqlc, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		sqlc = strings.Replace(sqlc, "~", dirname, 1)
		return sqlc + Extension()
	}

	if !strings.Contains(sqlc, string(filepath.Separator)) {
		return sqlc + Extension()
	}

	var pathError error
	sqlc, pathError = filepath.Rel(basepath, sqlc)
	if pathError != nil {
		Fatal(pathError)
	}

	return sqlc + Extension()
}
