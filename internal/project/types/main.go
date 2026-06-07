package main

import (
	"github.com/razshare/frizzante/internal/project/lib/core/types"
	"github.com/razshare/frizzante/internal/project/lib/routes/todos"
)

func main() {
	_ = types.Generate[todos.Props]()
}
