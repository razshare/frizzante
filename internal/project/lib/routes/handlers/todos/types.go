package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/types"
	session "github.com/razshare/frizzante/internal/project/lib/session/memory"
)

func init() {
	_ = types.Generate[Props]()
}

type Props struct {
	Todos []session.Todo `json:"todos"`
	Error string         `json:"error"`
}
