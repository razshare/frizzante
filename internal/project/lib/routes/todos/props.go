package todos

import (
	"github.com/razshare/frizzante/v2/internal/project/lib/schema"
)

type Props struct {
	Items []schema.Todo `json:"items"`
	Error string        `json:"error"`
}
