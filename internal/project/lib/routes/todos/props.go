package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/types"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func init() {
	_ = types.Generate[PropsForTodos]()
}

type PropsForTodos struct {
	Error string          `json:"error"`
	Items []sessions.Todo `json:"items"`
}
