//go:build !prod

package main

import (
	"log"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/dev/types"
	"github.com/razshare/frizzante/internal/project/lib/routes/todos"
)

func init() {
	directoryName := filepath.Join(".gen", "types")
	if err := types.Generate[todos.Props](directoryName); err != nil {
		log.Fatal(err)
		return
	}
}
