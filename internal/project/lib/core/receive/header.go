package receive

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// Header reads a header field and returns the value.
func Header(client *clients.Client, key string) string {
	return client.Request.Header.Get(key)
}

// ContentType reads the Content-Type header field and returns the value.
func ContentType(client *clients.Client) string {
	return client.Request.Header.Get("Content-Type")
}

// Accept reads if the Accept header entries and returns the values.
func Accept(client *clients.Client) string {
	return client.Request.Header.Get("Accept")
}
