package main

import (
	"embed"
	"github.com/razshare/frizzante/on"
)

//go:embed clilogo.txt
//go:embed database.sqlite
//go:embed version
//go:embed template/lib
//go:embed template/app/frizzante
//go:embed sqlc.yaml
//go:embed queries.sql
//go:embed schema.sql
var efs embed.FS

func main() {
	on.Start(efs)
}
