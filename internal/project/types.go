//go:build types

package main

import (
	"github.com/razshare/frizzante/internal/project/lib/core/types"
	"github.com/razshare/frizzante/internal/project/lib/routes/handlers/todos"
)

func main() {
	types.Generate[todos.Props]()
}
