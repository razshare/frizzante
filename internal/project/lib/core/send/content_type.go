package send

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// ContentType sets the Content-Type header field.
func ContentType(http *scopes.Http, ctype string) {
	Header(http, "Content-Type", ctype)
}
