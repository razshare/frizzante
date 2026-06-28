package logs

import (
	"fmt"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

func Infof(http *scopes.Http, format string, args ...any) {
	Info(http, fmt.Sprintf(format, args...))
}
