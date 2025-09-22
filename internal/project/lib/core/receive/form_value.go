package receive

import "github.com/razshare/frizzante/internal/project/lib/core/client"

// FormValue reads the first form value associated with the given key and returns it.
func FormValue(client *client.Client, key string) string {
	if !client.Parsed {
		Parse(client)
	}

	return client.Request.Form.Get(key)
}
