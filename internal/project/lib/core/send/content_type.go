package send

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// ContentType sets the Content-Type header field.
func ContentType(client *clients.Client, ctype string) {
	Header(client, "Content-Type", ctype)
}
