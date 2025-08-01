package archives

import (
	"github.com/razshare/frizzante/locks"
)

type DiskArchive struct {
	Name string
	Lock *locks.Lock
}

type Archive interface {
	Get(domain string, key string) ([]byte, error)
	Set(domain string, key string, value []byte) error
	Has(domain string, key string) (bool, error)
	Remove(domain string, key string) error
	HasDomain(domain string) (bool, error)
	RemoveDomain(domain string) error
}
