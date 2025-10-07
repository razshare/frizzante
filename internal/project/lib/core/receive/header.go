package receive

import _client "github.com/razshare/frizzante/internal/project/lib/core/client"

// Header reads a header field and returns the value.
func Header(client *_client.Client, key string) string {
	return client.Request.Header.Get(key)
}

// ContentType reads the Content-Type header field and returns the value.
func ContentType(client *_client.Client) string {
	return client.Request.Header.Get("Content-Type")
}

// Accept reads if the Accept header entries and returns the values.
func Accept(client *_client.Client) string {
	return client.Request.Header.Get("Accept")
}
