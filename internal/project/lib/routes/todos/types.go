package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/types"
	"github.com/razshare/frizzante/internal/project/lib/memory/sessions"
)

func init() {
	_ = types.Generate[Props]()
}

type Props struct {
	Error string          `json:"error"`
	Items []sessions.Todo `json:"items"`
}

type ToggleForm struct {
	Index int `form:"index"`
	Value int `form:"value"`
}

type AddForm struct {
	Description string `form:"description"`
}

type RemoveForm struct {
	Index int `form:"index"`
}
