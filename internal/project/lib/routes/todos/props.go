package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/databases/sqlc"
)

type Props struct {
	Items []sqlc.Todo `json:"items"`
	Error string      `json:"error"`
}
