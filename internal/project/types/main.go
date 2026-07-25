package main

import (
	"github.com/razshare/frizzante/v2/internal/project/lib/core/types"
	"github.com/razshare/frizzante/v2/internal/project/lib/routes/todos"
)

func main() {
	_ = types.Clear()
	_ = types.Generate[todos.Props]()
}
