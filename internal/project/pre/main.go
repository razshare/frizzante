package main

import (
	"log"

	"github.com/razshare/frizzante/internal/project/lib/core/types"
	"github.com/razshare/frizzante/internal/project/lib/routes/todos"
)

var errors = []error{
	// add your shared types here
	types.Generate[todos.Props](),
}

func main() {
	for _, err := range errors {
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("types generates in .gen/types")
}
