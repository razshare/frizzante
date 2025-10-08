package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/types"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func init() {
	_ = types.Generate[Props]()
}

type Props struct {
	Items []sessions.Todo `json:"items"`
	Error string          `json:"error"`
}
