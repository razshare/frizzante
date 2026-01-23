package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

type Props struct {
	Error string          `json:"error"`
	Items []sessions.Todo `json:"items"`
}
