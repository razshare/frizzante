package libarchive

import (
	"github.com/razshare/frizzante/libroad"
)

type Archive interface {
	Get(domain string, key string) ([]byte, error)
	Set(domain string, key string, value []byte) error
	Has(domain string, key string) (bool, error)
	Remove(domain string, key string) error
	HasDomain(domain string) (bool, error)
	RemoveDomain(domain string) error
}

type DiskArchive struct {
	Name string
	Road *libroad.Road
}
