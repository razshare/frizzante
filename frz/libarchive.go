package frz

import (
	"errors"
	"github.com/razshare/frizzante/fs"
	"os"
	"path/filepath"
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
	name string
	road *Road
}

// NewDiskArchive create a new archive backed by the file system.
func NewDiskArchive() *DiskArchive {
	return &DiskArchive{
		name: "archive",
		road: NewRoad(),
	}
}

// WithName sets the name of the archive and thus its root directory.
func (archive *DiskArchive) WithName(name string) *DiskArchive {
	archive.name = name
	return archive
}

// Set sets a value to the archive based on the domain and key.
func (archive *DiskArchive) Set(domain string, key string, value []byte) error {
	if "" == archive.name {
		return errors.New("disk archive name is blank")
	}
	lane := archive.road.WithLane(domain, key)
	lane.Lock()
	defer lane.Unlock()
	directoryName := filepath.Join(archive.name, domain)
	if !fs.FileExists(directoryName) {
		mkdirError := os.MkdirAll(directoryName, os.ModePerm)
		if nil != mkdirError {
			return mkdirError
		}
	}
	fileName := filepath.Join(directoryName, key)
	writeError := os.WriteFile(fileName, value, os.ModePerm)
	if nil != writeError {
		return writeError
	}
	return nil
}

// Get gets a value from the archive based on domain and key.
func (archive *DiskArchive) Get(domain string, key string) ([]byte, error) {
	if "" == archive.name {
		return nil, errors.New("disk archive name is blank")
	}
	lane := archive.road.WithLane(domain, key)
	lane.Lock()
	defer lane.Unlock()
	fileName := filepath.Join(archive.name, domain, key)
	value, readError := os.ReadFile(fileName)
	if nil != readError {
		return nil, readError
	}
	return value, nil
}

// Has checks if the archive has a value based on a domain and key.
func (archive *DiskArchive) Has(domain string, key string) (bool, error) {
	if "" == archive.name {
		return false, errors.New("disk archive name is blank")
	}
	lane := archive.road.WithLane(domain, key)
	lane.Lock()
	defer lane.Unlock()
	fileName := filepath.Join(archive.name, domain, key)
	return fs.FileExists(fileName), nil
}

// Remove removes a value from the archive based on a domain and key.
func (archive *DiskArchive) Remove(domain string, key string) error {
	if "" == archive.name {
		return errors.New("disk archive name is blank")
	}
	lane := archive.road.WithLane(domain, key)
	lane.Lock()
	defer lane.Unlock()
	fileName := filepath.Join(archive.name, domain, key)
	removeError := os.Remove(fileName)
	if nil != removeError {
		return removeError
	}
	return nil
}

// HasDomain checks if the archive has a domain.
func (archive *DiskArchive) HasDomain(domain string) (bool, error) {
	if "" == archive.name {
		return false, errors.New("disk archive name is blank")
	}
	lane := archive.road.WithLane(domain)
	lane.Lock()
	defer lane.Unlock()
	return fs.FileExists(filepath.Join(archive.name, domain)), nil
}

// RemoveDomain removes a domain from the archive.
func (archive *DiskArchive) RemoveDomain(domain string) error {
	if "" == archive.name {
		return errors.New("disk archive name is blank")
	}
	lane := archive.road.WithLane(domain)
	lane.Lock()
	defer lane.Unlock()
	directoryName := filepath.Join(archive.name, domain)
	removeError := os.RemoveAll(directoryName)
	if nil != removeError {
		return removeError
	}
	return nil
}
