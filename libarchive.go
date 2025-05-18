package frizzante

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type ArchiveBuilder = func(archive *Archive)

type Archive struct {
	name string

	get    func(domain string, key string) []byte
	set    func(domain string, key string, value []byte)
	has    func(domain string, key string) bool
	remove func(domain string, key string)

	domainExists func(domain string) bool
	removeDomain func(domain string)

	notifier *Notifier
}

// Get reads the combination of domain and key from the archive.
func (archive *Archive) Get(domain string, key string) []byte {
	return archive.get(domain, key)
}

// Set writes to the combination of domain and key int the archive.
func (archive *Archive) Set(domain string, key string, value []byte) {
	archive.set(domain, key, value)
}

// Has checks if the combination of domain and key exists in the archive.
func (archive *Archive) Has(domain string, key string) bool {
	return archive.has(domain, key)
}

// Remove removes the combination of domain and key from the archive.
func (archive *Archive) Remove(domain string, key string) {
	archive.remove(domain, key)
}

// DomainExists checks if the archive has a domain.
func (archive *Archive) DomainExists(domain string) bool {
	return archive.domainExists(domain)
}

// RemoveDomain removes a domain from the archive.
func (archive *Archive) RemoveDomain(domain string) {
	archive.removeDomain(domain)
}

// WithName sets the name of the archive.
func (archive *Archive) WithName(name string) {
	archive.name = name
}

// WithKeyGetter sets the reader, which reads a value from the archive
// given its domain and key.
func (archive *Archive) WithKeyGetter(getter func(domain string, key string) []byte) {
	archive.get = func(domain string, key string) []byte {
		if !KeyIsSafe(domain) {
			archive.notifier.SendError(
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", archive.name, domain))
			return nil
		}

		if !KeyIsSafe(key) {
			archive.notifier.SendError(
				fmt.Errorf("archive `%s` has rejected key `%s` because it looks malicious", archive.name, key))
			return nil
		}

		return getter(domain, key)
	}
}

// WithKeySetter sets the writer,
// which writes a value into the archive at a combination of
// domain and key.
func (archive *Archive) WithKeySetter(setter func(domain string, key string, value []byte)) {
	archive.set = func(domain string, key string, value []byte) {
		if !KeyIsSafe(domain) {
			archive.notifier.SendError(
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", archive.name, domain))
			return
		}

		if !KeyIsSafe(key) {
			archive.notifier.SendError(
				fmt.Errorf("archive `%s` has rejected key `%s` because it looks malicious", archive.name, key))
			return
		}

		setter(domain, key, value)
	}
}

// WithDomainRemover sets the remover,
// which removes a domain from the archive.
func (archive *Archive) WithDomainRemover(domainRemover func(domain string)) {
	archive.removeDomain = func(domain string) {
		if !KeyIsSafe(domain) {
			archive.notifier.SendError(
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", archive.name, domain))
			return
		}

		domainRemover(domain)
	}
}

// WithKeyRemover sets the remover,
// which removes a combination of domain and key from the archive.
func (archive *Archive) WithKeyRemover(remover func(domain string, key string)) {
	archive.remove = func(domain string, key string) {
		if !KeyIsSafe(domain) {
			archive.notifier.SendError(
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", archive.name, domain))
			return
		}

		if !KeyIsSafe(key) {
			archive.notifier.SendError(
				fmt.Errorf("archive `%s` has rejected key `%s` because it looks malicious", archive.name, key))
			return
		}

		remover(domain, key)
	}
}

// WithDomainChecker sets the domain checker,
// which checks if a domain exists in the archive.
func (archive *Archive) WithDomainChecker(checker func(domain string) (exists bool)) {
	archive.domainExists = func(domain string) (exists bool) {
		if !KeyIsSafe(domain) {
			archive.notifier.SendError(
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", archive.name, domain))
			return
		}

		return checker(domain)
	}
}

// WithKeyChecker sets the domain and key checker,
// which checks if a combination of domain and key exists in the archive.
func (archive *Archive) WithKeyChecker(checker func(domain string, key string) (exists bool)) {
	archive.has = func(domain string, key string) (exists bool) {
		if !KeyIsSafe(domain) {
			archive.notifier.SendError(
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", archive.name, domain))
			return
		}

		if !KeyIsSafe(key) {
			archive.notifier.SendError(
				fmt.Errorf("archive `%s` has rejected key `%s` because it looks malicious", archive.name, key))
			return
		}

		return checker(domain, key)
	}
}

// WithNotifier sets the notifier.
func (archive *Archive) WithNotifier(notifier *Notifier) {
	archive.notifier = notifier
}

// NewArchive creates an archive.
//
// The alphabet of the archive is composed of the english alphabet (upper case letters and lowercase letters),
// digits from 0 to 9, the "_" (underscore) character and the "-" (dash) character.
//
// Any domain or key received by the archive that is not in scope of the alphabet will be rejected automatically,
// regardless of the custom implementations of the getter, setter, remover and checker functions.
func NewArchive(builder ArchiveBuilder) *Archive {
	archiveName, archiveNameError := filepath.Abs("archive")
	if nil != archiveNameError {
		log.Fatal(archiveNameError)
	}
	archive := &Archive{
		name:     archiveName,
		notifier: NewNotifier(),
	}

	builder(archive)

	return archive
}

// NewArchiveOnDisk creates an archive that uses the local file system as a backend.
func NewArchiveOnDisk(name string, cacheTtl time.Duration) *Archive {
	keys := NewCache[[]byte]()
	domains := NewCache[bool]()
	road := NewRoad()

	return NewArchive(func(a *Archive) {
		a.WithName(name)
		a.WithKeyGetter(func(domain string, key string) []byte {
			lane := road.WithLane(domain, key)
			<-lane
			fileName := filepath.Join(a.name, domain, key)
			if keys.IsNotExpired(fileName) {
				value := keys.Get(fileName)
				lane <- 0
				return value
			}
			value, readError := os.ReadFile(fileName)
			if nil != readError {
				a.notifier.SendError(readError)
				lane <- 0
				return nil
			}
			keys.Set(cacheTtl, fileName, value)
			lane <- 0
			return value
		})

		a.WithKeySetter(func(domain string, key string, value []byte) {
			lane := road.WithLane(domain, key)
			<-lane
			directoryName := filepath.Join(a.name, domain)
			if !fileExists(directoryName) {
				mkdirError := os.MkdirAll(directoryName, os.ModePerm)
				if nil != mkdirError {
					a.notifier.SendError(mkdirError)
					lane <- 0
					return
				}
			}
			fileName := filepath.Join(directoryName, key)
			writeError := os.WriteFile(fileName, value, os.ModePerm)
			if nil != writeError {
				a.notifier.SendError(writeError)
			}
			keys.Set(cacheTtl, fileName, value)
			lane <- 0
		})

		a.WithDomainChecker(func(domain string) bool {
			lane := road.WithLane(domain)
			<-lane
			directoryName := filepath.Join(a.name, domain)
			if domains.IsNotExpired(directoryName) {
				value := domains.Get(directoryName)
				lane <- 0
				return value
			}
			ok := fileExists(directoryName)
			domains.Set(cacheTtl, directoryName, ok)
			lane <- 0
			return ok
		})

		a.WithKeyChecker(func(domain string, key string) bool {
			lane := road.WithLane(domain, key)
			<-lane
			fileName := filepath.Join(a.name, domain, key)
			if domains.IsNotExpired(fileName) {
				value := domains.Get(fileName)
				lane <- 0
				return value
			}
			ok := fileExists(fileName)
			domains.Set(cacheTtl, fileName, ok)
			lane <- 0
			return ok
		})

		a.WithDomainRemover(func(domain string) {
			lane := road.WithLane(domain)
			<-lane
			directoryName := filepath.Join(a.name, domain)
			removeError := os.RemoveAll(directoryName)
			if nil != removeError {
				a.notifier.SendError(removeError)
			}
			keys.Remove(directoryName)
			lane <- 0
		})

		a.WithKeyRemover(func(domain string, key string) {
			lane := road.WithLane(domain, key)
			<-lane
			fileName := filepath.Join(a.name, domain, key)
			removeError := os.Remove(fileName)
			if nil != removeError {
				a.notifier.SendError(removeError)
				lane <- 0
				return
			}
			keys.Remove(fileName)
			lane <- 0
		})
	})
}
