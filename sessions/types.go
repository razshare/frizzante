package sessions

import (
	"github.com/razshare/frizzante/connections"
)

type Session[T any] struct {
	// Get gets a value from the archive based on domain and key.
	Get func(domain string, key string) ([]byte, error)
	// Set sets a value to the archive based on the domain and key.
	Set func(domain string, key string, value []byte) error
	// Has checks if the archive has a value based on a domain and key.
	Has func(domain string, key string) (bool, error)
	// Remove removes a value from the archive based on a domain and key.
	Remove func(domain string, key string) error
	// HasDomain checks if the archive has a domain.
	HasDomain func(domain string) (bool, error)
	// RemoveDomain removes a domain from the archive.
	RemoveDomain func(domain string) error
	Connection   *connections.Connection
	State        *T
}
