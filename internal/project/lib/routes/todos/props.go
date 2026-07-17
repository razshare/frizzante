package todos

import "github.com/razshare/frizzante/internal/project/lib/databases/schema"

type Props struct {
	Items []schema.Todo `json:"items"`
	Error string        `json:"error"`
}
