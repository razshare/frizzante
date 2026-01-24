package main

import (
	"log"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/dev/types"
	"github.com/razshare/frizzante/internal/project/lib/routes/todos"
)

var directoryName = filepath.Join(".gen", "types")

func main() {
	if err := types.Generate[todos.Props](directoryName); err != nil {
		log.Fatal(err)
		return
	}
}
