package receive

import "main/lib/core/client"

// Header reads a header field and returns the value.
//
// Compatible with web sockets.
func Header(client *client.Client, key string) string {
	return client.Request.Header.Get(key)
}

// ContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func ContentType(client *client.Client) string {
	return client.Request.Header.Get("Content-Type")
}

// Accept reads if the Accept header entries and returns the values.
func Accept(client *client.Client) string {
	return client.Request.Header.Get("Accept")
}
