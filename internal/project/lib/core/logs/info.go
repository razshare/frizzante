package logs

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

func Info(http *scopes.Http, args ...any) {
	length := len(args)
	entries := make([]string, len(args))
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	http.InfoLog.Println(strings.Join(entries, ""))
}
