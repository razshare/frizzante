package main

import (
	"log"

	"github.com/razshare/frizzante/internal/project/lib/core/dev/types"
	"github.com/razshare/frizzante/internal/project/lib/routes/todos"
)

func main() {
	for _, err := range []error{
		// add your shared types here
		types.Generate[todos.Props](),
	} {
		if err != nil {
			log.Fatal(err)
		}
	}
}
