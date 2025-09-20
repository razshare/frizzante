package todos

import session "github.com/razshare/frizzante/internal/project/lib/session/memory"

type Props struct {
	Todos []session.Todo `json:"todos"`
	Error string         `json:"error"`
}
