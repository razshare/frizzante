package main

import (
	"log"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/types"
	"github.com/razshare/frizzante/internal/project/lib/routes/todos"
)

var directoryName = filepath.Join(".gen", "types")
var errors = []error{
	// add your shared types here
	types.Generate[todos.Props](directoryName),
}

func main() {
	for _, err := range errors {
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("types generates in %s", directoryName)
}
