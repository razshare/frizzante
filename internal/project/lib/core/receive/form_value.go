package receive

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormValue reads the first form value associated with the given key and returns it.
func FormValue(client *client.Client, key string) string {
	if client.Request.Form == nil {
		if err := client.Request.ParseMultipartForm(MaxFormSize); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return ""
		}
	}
	if vs := client.Request.Form[key]; len(vs) > 0 {
		return vs[0]
	}
	return ""
}
