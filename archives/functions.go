package archives

import (
	"errors"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/locks"
	"os"
	"path/filepath"
)

// New creates a new archive backed by the file system.
func New(name string) Archive {
	return &DiskArchive{
		Name: name,
		Lock: locks.New(),
	}
}

// Set sets a value to the archive based on the domain and key.
func (archive *DiskArchive) Set(domain string, key string, value []byte) error {
	if "" == archive.Name {
		return errors.New("disk archive name is blank")
	}

	mutex := archive.Lock.Acquire(domain, key)
	mutex.Lock()
	defer mutex.Unlock()
	directoryName := filepath.Join(archive.Name, domain)
	if !files.IsDirectory(directoryName) {
		mkdirError := os.MkdirAll(directoryName, os.ModePerm)
		if mkdirError != nil {
			return mkdirError
		}
	}
	fileName := filepath.Join(directoryName, key)
	writeError := os.WriteFile(fileName, value, os.ModePerm)
	if writeError != nil {
		return writeError
	}
	return nil
}

// Get gets a value from the archive based on domain and key.
func (archive *DiskArchive) Get(domain string, key string) ([]byte, error) {
	if "" == archive.Name {
		return nil, errors.New("disk archive name is blank")
	}
	mutex := archive.Lock.Acquire(domain, key)
	mutex.Lock()
	defer mutex.Unlock()
	fileName := filepath.Join(archive.Name, domain, key)
	value, readError := os.ReadFile(fileName)
	if readError != nil {
		return nil, readError
	}
	return value, nil
}

// Has checks if the archive has a value based on a domain and key.
func (archive *DiskArchive) Has(domain string, key string) (bool, error) {
	if "" == archive.Name {
		return false, errors.New("disk archive name is blank")
	}
	mutex := archive.Lock.Acquire(domain, key)
	mutex.Lock()
	defer mutex.Unlock()
	fileName := filepath.Join(archive.Name, domain, key)
	return files.IsFile(fileName), nil
}

// Remove removes a value from the archive based on a domain and key.
func (archive *DiskArchive) Remove(domain string, key string) error {
	if "" == archive.Name {
		return errors.New("disk archive name is blank")
	}

	mutex := archive.Lock.Acquire(domain, key)
	mutex.Lock()
	defer mutex.Unlock()

	fileName := filepath.Join(archive.Name, domain, key)
	removeError := os.Remove(fileName)
	if removeError != nil {
		return removeError
	}
	return nil
}

// HasDomain checks if the archive has a domain.
func (archive *DiskArchive) HasDomain(domain string) (bool, error) {
	if "" == archive.Name {
		return false, errors.New("disk archive name is blank")
	}
	mutex := archive.Lock.Acquire(domain)
	mutex.Lock()
	defer mutex.Unlock()
	return files.IsDirectory(filepath.Join(archive.Name, domain)), nil
}

// RemoveDomain removes a domain from the archive.
func (archive *DiskArchive) RemoveDomain(domain string) error {
	if "" == archive.Name {
		return errors.New("disk archive name is blank")
	}
	mutex := archive.Lock.Acquire(domain)
	mutex.Lock()
	defer mutex.Unlock()
	directoryName := filepath.Join(archive.Name, domain)
	removeError := os.RemoveAll(directoryName)
	if removeError != nil {
		return removeError
	}
	return nil
}
