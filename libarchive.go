package frizzante

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

type ArchiveBuilder = func(archive *ArchiveInterface)

type ArchiveInterface interface {
	Get(domain string, key string) []byte
	Set(domain string, key string, value []byte)
	Has(domain string, key string) bool
	Remove(domain string, key string)
	HasDomain(domain string) bool
	RemoveDomain(domain string)
}

type DiskArchive struct {
	name     string
	ttl      time.Duration
	keys     *Cache
	domains  *Cache
	road     *Road
	notifier Notifier
}

func NewDiskArchive(name string, ttl time.Duration, notifier Notifier) *DiskArchive {
	return &DiskArchive{
		name:     name,
		ttl:      ttl,
		notifier: notifier,
		keys:     NewCache(),
		domains:  NewCache(),
		road:     NewRoad(),
	}
}

func (archive *DiskArchive) Get(domain string, key string) []byte {
	if "" == archive.name {
		archive.notifier.SendError(errors.New("disk archive name is blank"))
		return make([]byte, 0)
	}
	lane := archive.road.WithLane(domain, key)
	<-lane
	fileName := filepath.Join(archive.name, domain, key)
	if archive.keys.IsNotExpired(fileName) {
		value := archive.keys.Get(fileName).([]byte)
		lane <- 0
		return value
	}
	value, readError := os.ReadFile(fileName)
	if nil != readError {
		lane <- 0
		archive.notifier.SendError(readError)
		return make([]byte, 0)
	}
	archive.keys.Set(archive.ttl, fileName, value)
	lane <- 0
	return value
}

func (archive *DiskArchive) Set(domain string, key string, value []byte) {
	if "" == archive.name {
		archive.notifier.SendError(errors.New("disk archive name is blank"))
		return
	}
	lane := archive.road.WithLane(domain, key)
	<-lane
	directoryName := filepath.Join(archive.name, domain)
	if !FileExists(directoryName) {
		mkdirError := os.MkdirAll(directoryName, os.ModePerm)
		if nil != mkdirError {
			lane <- 0
			archive.notifier.SendError(mkdirError)
			return
		}
	}
	fileName := filepath.Join(directoryName, key)
	writeError := os.WriteFile(fileName, value, os.ModePerm)
	if nil != writeError {
		lane <- 0
		archive.notifier.SendError(writeError)
		return
	}
	archive.keys.Set(archive.ttl, fileName, value)
	lane <- 0
}

func (archive *DiskArchive) Has(domain string, key string) bool {
	if "" == archive.name {
		archive.notifier.SendError(errors.New("disk archive name is blank"))
		return false
	}
	lane := archive.road.WithLane(domain, key)
	<-lane
	fileName := filepath.Join(archive.name, domain, key)
	if archive.domains.IsNotExpired(fileName) {
		value := archive.domains.Get(fileName).(bool)
		lane <- 0
		return value
	}
	ok := FileExists(fileName)
	archive.domains.Set(archive.ttl, fileName, ok)
	lane <- 0
	return ok
}

func (archive *DiskArchive) Remove(domain string, key string) {
	if "" == archive.name {
		archive.notifier.SendError(errors.New("disk archive name is blank"))
		return
	}
	lane := archive.road.WithLane(domain, key)
	<-lane
	fileName := filepath.Join(archive.name, domain, key)
	removeError := os.Remove(fileName)
	if nil != removeError {
		lane <- 0
		archive.notifier.SendError(removeError)
		return
	}
	archive.keys.Remove(fileName)
	lane <- 0
}

func (archive *DiskArchive) HasDomain(domain string) bool {
	if "" == archive.name {
		archive.notifier.SendError(errors.New("disk archive name is blank"))
		return false
	}
	lane := archive.road.WithLane(domain)
	<-lane
	directoryName := filepath.Join(archive.name, domain)
	if archive.domains.IsNotExpired(directoryName) {
		value := archive.domains.Get(directoryName).(bool)
		lane <- 0
		return value
	}
	ok := FileExists(directoryName)
	archive.domains.Set(archive.ttl, directoryName, ok)
	lane <- 0
	return ok
}

func (archive *DiskArchive) RemoveDomain(domain string) {
	if "" == archive.name {
		archive.notifier.SendError(errors.New("disk archive name is blank"))
		return
	}
	lane := archive.road.WithLane(domain)
	<-lane
	directoryName := filepath.Join(archive.name, domain)
	removeError := os.RemoveAll(directoryName)
	if nil != removeError {
		archive.notifier.SendError(removeError)
	}
	archive.keys.Remove(directoryName)
	lane <- 0
}
