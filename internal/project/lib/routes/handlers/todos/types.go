package todos

import session "github.com/razshare/frizzante/internal/project/lib/session/memory"

type Props struct {
	Todos []session.Todo
	Error string
}
