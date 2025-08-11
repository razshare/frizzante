package download

import (
	"embed"
	"github.com/razshare/frizzante/cli/platform"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func Sqlc(efs embed.FS) {
	dn := filepath.Join(".gen", "sqlc")
	p := platform.Find()

	var url string

	if p == platform.DarwinArm64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_arm64.zip"
	} else if p == platform.DarwinAmd64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_amd64.zip"
	} else if p == platform.LinuxArm64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_arm64.zip"
	} else if p == platform.LinuxAmd64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_amd64.zip"
	} else if p == platform.WindowsArm64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
	} else if p == platform.WindowsAmd64 {
		url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
	}

	Install("sqlc", url, dn)

	schema := true
	queries := true
	yaml := true

	if files.IsFile("schema.sql") {
		schema = confirm.Send(true, "File `schema.sql` already exists, would you like to overwrite it?")
	}

	if schema {
		err := os.WriteFile("schema.sql", make([]byte, 0), os.ModePerm)
		if err != nil {
			messages.Fatal(err)
		}
		messages.Success("schema.sql created")
	}

	if files.IsFile("queries.sql") {
		queries = confirm.Send(true, "File `queries.sql` already exists, would you like to overwrite it?")
	}

	if queries {
		err := os.WriteFile("queries.sql", make([]byte, 0), os.ModePerm)
		if err != nil {
			messages.Fatal(err)
		}
		messages.Success("queries.sql created")
	}

	if files.IsFile("sqlc.yaml") {
		yaml = confirm.Send(true, "File `sqlc.yaml` already exists, would you like to overwrite it?")
	}

	if yaml {
		data, err := efs.ReadFile("sqlc.yaml")
		if err != nil {
			messages.Fatal(err)
		}

		err = os.WriteFile("sqlc.yaml", data, os.ModePerm)
		if err != nil {
			messages.Fatal(err)
		}
		messages.Success("sqlc.yaml created")
	}

	//if !files.IsFile("database.sqlite") {
	//	messages.Info("frizzante's sqlc configuration uses sqlite by default")
	//	if confirm.Sendf(true, "would you like to create an empty sqlite database?") {
	//		OnSqliteDatabase(efs)
	//	}
	//}

	return
}
