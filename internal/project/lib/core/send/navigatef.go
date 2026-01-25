package send

import (
	"fmt"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// Navigatef redirects the request to a location with status 302.
func Navigatef(client *clients.Client, format string, vars ...any) {
	Redirect(client, fmt.Sprintf(format, vars...), 302)
	Message(client, "")
}
