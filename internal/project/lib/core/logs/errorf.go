package logs

import (
	"fmt"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

func Errorf(http *scopes.Http, format string, args ...any) {
	Error(http, fmt.Sprintf(format, args...))
}
