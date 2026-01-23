package main

import (
	"log"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/types"
	"github.com/razshare/frizzante/internal/project/lib/routes/todos"
)

var directory = filepath.Join(".gen", "types")

func main() {
	if err := types.Generate[todos.Props](directory); err != nil {
		log.Fatal(err)
	}
}
