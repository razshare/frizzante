package send

import (
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Headers sends header fields.
func Headers(http *scopes.Http, fields map[string]string) {
	if http.Locked {
		logs.Errorf(
			http,
			"send.Headers: headers are locked, cannot set headers\n%s",
			stack.Trace(),
		)
		return
	}
	for key, value := range fields {
		http.Writer.Header().Set(key, value)
	}
}
