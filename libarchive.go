package frizzante

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
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
	alphabet []rune
}

// ArchiveGet reads the combination of domain and key from the archive.
func ArchiveGet(self *Archive, domain string, key string) []byte {
	return self.get(domain, key)
}

// ArchiveSet writes to the combination of domain and key int the archive.
func ArchiveSet(self *Archive, domain string, key string, value []byte) {
	self.set(domain, key, value)
}

// ArchiveHas checks if the combination of domain and key exists in the archive.
func ArchiveHas(self *Archive, domain string, key string) bool {
	return self.has(domain, key)
}

// ArchiveRemove removes the combination of domain and key from the archive.
func ArchiveRemove(self *Archive, domain string, key string) {
	self.remove(domain, key)
}

// ArchiveDomainExists checks if the archive has a domain.
func ArchiveDomainExists(self *Archive, domain string) bool {
	return self.domainExists(domain)
}

// ArchiveRemoveDomain removes a domain from the archive.
func ArchiveRemoveDomain(self *Archive, domain string) {
	self.removeDomain(domain)
}

// ArchiveAcceptsDomain checks if a domain is accepted by the alphabet of the archive.
//
// Generally speaking alphabets should not accept runes like "/", "\", ".", ".." and so on.
//
// This is so that "sneaky" or maliciously constructed domains injected by clients
// can't change directories.
//
// You don't have to use ArchiveAcceptsDomain to manually check for these malicious strings,
// this check is executed automatically by the archive internally, regardless of your implementation.
//
// This function is public only as a quality of life improvement, so that if you would like to
// detect such strings before crashing the archive, you have a way to do so.
func ArchiveAcceptsDomain(self *Archive, domain string) bool {
	for _, char := range domain {
		if !slices.Contains(self.alphabet, char) {
			return false
		}
	}

	return true
}

// ArchiveAcceptsKey checks if a key is accepted by the alphabet of the archive.
//
// Generally speaking alphabets should not accept runes like "/", "\", ".", ".." and so on.
//
// This is so that "sneaky" or maliciously constructed keys injected by clients
// can't change directories.
//
// You don't have to use ArchiveAcceptsDomain to manually check for these malicious strings,
// this check is executed automatically by the archive internally, regardless of your implementation.
//
// This function is public only as a quality of life improvement, so that if you would like to
// detect such strings before crashing the archive, you have a way to do so.
func ArchiveAcceptsKey(self *Archive, key string) bool {
	for _, char := range key {
		if !slices.Contains(self.alphabet, char) {
			return false
		}
	}

	return true
}

// ArchiveWithName sets the name of the archive.
func ArchiveWithName(self *Archive, name string) {
	self.name = name
}

// ArchiveWithKeyGetter sets the reader, which reads a value from the archive
// given its domain and key.
func ArchiveWithKeyGetter(self *Archive, getter func(domain string, key string) []byte) {
	self.get = func(domain string, key string) []byte {
		if !ArchiveAcceptsDomain(self, domain) {
			NotifierSendError(self.notifier,
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", self.name, domain))
			return nil
		}

		if !ArchiveAcceptsKey(self, key) {
			NotifierSendError(self.notifier,
				fmt.Errorf("archive `%s` has rejected key `%s` because it looks malicious", self.name, key))
			return nil
		}

		return getter(domain, key)
	}
}

// ArchiveWithKeySetter sets the writer,
// which writes a value into the archive at a combination of
// domain and key.
func ArchiveWithKeySetter(self *Archive, setter func(domain string, key string, value []byte)) {
	self.set = func(domain string, key string, value []byte) {
		if !ArchiveAcceptsDomain(self, domain) {
			NotifierSendError(self.notifier,
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", self.name, domain))
			return
		}

		if !ArchiveAcceptsKey(self, key) {
			NotifierSendError(self.notifier,
				fmt.Errorf("archive `%s` has rejected key `%s` because it looks malicious", self.name, key))
			return
		}

		setter(domain, key, value)
	}
}

// ArchiveWithDomainRemover sets the remover,
// which removes a domain from the archive.
func ArchiveWithDomainRemover(self *Archive, domainRemover func(domain string)) {
	self.removeDomain = func(domain string) {
		if !ArchiveAcceptsDomain(self, domain) {
			NotifierSendError(self.notifier,
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", self.name, domain))
			return
		}

		domainRemover(domain)
	}
}

// ArchiveWithKeyRemover sets the remover,
// which removes a combination of domain and key from the archive.
func ArchiveWithKeyRemover(self *Archive, remover func(domain string, key string)) {
	self.remove = func(domain string, key string) {
		if !ArchiveAcceptsDomain(self, domain) {
			NotifierSendError(self.notifier,
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", self.name, domain))
			return
		}

		if !ArchiveAcceptsKey(self, key) {
			NotifierSendError(self.notifier,
				fmt.Errorf("archive `%s` has rejected key `%s` because it looks malicious", self.name, key))
			return
		}

		remover(domain, key)
	}
}

// ArchiveWithDomainChecker sets the domain checker,
// which checks if a domain exists in the archive.
func ArchiveWithDomainChecker(self *Archive, checker func(domain string) (exists bool)) {
	self.domainExists = func(domain string) (exists bool) {
		if !ArchiveAcceptsDomain(self, domain) {
			NotifierSendError(self.notifier,
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", self.name, domain))
			return
		}

		return checker(domain)
	}
}

// ArchiveWithKeyChecker sets the domain and key checker,
// which checks if a combination of domain and key exists in the archive.
func ArchiveWithKeyChecker(self *Archive, checker func(domain string, key string) (exists bool)) {
	self.has = func(domain string, key string) (exists bool) {
		if !ArchiveAcceptsDomain(self, domain) {
			NotifierSendError(self.notifier,
				fmt.Errorf("archive `%s` has rejected domain `%s` because it looks malicious", self.name, domain))
			return
		}

		if !ArchiveAcceptsKey(self, key) {
			NotifierSendError(self.notifier,
				fmt.Errorf("archive `%s` has rejected key `%s` because it looks malicious", self.name, key))
			return
		}

		return checker(domain, key)
	}
}

// ArchiveWithNotifier sets the notifier.
func ArchiveWithNotifier(self *Archive, notifier *Notifier) {
	self.notifier = notifier
}

// ArchiveCreate creates an archive.
//
// The alphabet of the archive is composed of the english alphabet (upper case letters and lowercase letters),
// digits from 0 to 9, the "_" (underscore) character and the "-" (dash) character.
//
// Any domain or key received by the archive that is not in scope of the alphabet will be rejected automatically,
// regardless of the custom implementations of the getter, setter, remover and checker functions.
func ArchiveCreate(builder ArchiveBuilder) *Archive {
	archiveName, archiveNameError := filepath.Abs("archive")
	if nil != archiveNameError {
		log.Fatal(archiveNameError)
	}
	archive := &Archive{
		name: archiveName,
		alphabet: []rune{
			'A',
			'B',
			'C',
			'D',
			'E',
			'F',
			'G',
			'H',
			'I',
			'J',
			'K',
			'L',
			'M',
			'N',
			'O',
			'P',
			'Q',
			'R',
			'S',
			'T',
			'U',
			'V',
			'X',
			'Y',
			'Z',
			'a',
			'b',
			'c',
			'd',
			'e',
			'f',
			'g',
			'h',
			'i',
			'j',
			'k',
			'l',
			'm',
			'n',
			'o',
			'p',
			'q',
			'r',
			's',
			't',
			'u',
			'v',
			'x',
			'y',
			'z',
			'1',
			'2',
			'3',
			'4',
			'5',
			'6',
			'7',
			'8',
			'9',
			'0',
			'_',
			'-',
			'.',
		},
	}

	builder(archive)

	return archive
}

// ArchiveCreateOnDisk creates an archive that uses the local file system as a backend.
func ArchiveCreateOnDisk(name string) *Archive {
	return ArchiveCreate(func(archive *Archive) {
		ArchiveWithName(archive, name)
		ArchiveWithKeyGetter(archive, func(domain string, key string) []byte {
			fileName := filepath.Join(archive.name, domain, key)
			content, readError := os.ReadFile(fileName)
			if nil != readError && nil != archive.notifier {
				NotifierSendError(archive.notifier, readError)
				return nil
			}
			return content
		})

		ArchiveWithKeySetter(archive, func(domain string, key string, value []byte) {
			directoryName := filepath.Join(archive.name, domain)

			if !exists(directoryName) {
				mkdirError := os.MkdirAll(directoryName, os.ModePerm)
				if nil != mkdirError {
					NotifierSendError(archive.notifier, mkdirError)
					return
				}
			}

			fileName := filepath.Join(archive.name, domain, key)
			writeError := os.WriteFile(fileName, value, os.ModePerm)
			if nil != writeError {
				NotifierSendError(archive.notifier, writeError)
			}
		})

		ArchiveWithKeyChecker(archive, func(domain string, key string) bool {
			fileName := filepath.Join(archive.name, domain, key)
			return exists(fileName)
		})

		ArchiveWithKeyRemover(archive, func(domain string, key string) {
			fileName := filepath.Join(archive.name, domain, key)
			removeError := os.Remove(fileName)
			if nil != removeError {
				NotifierSendError(archive.notifier, removeError)
			}
		})

		ArchiveWithDomainChecker(archive, func(domain string) bool {
			fileName := filepath.Join(archive.name, domain)
			return exists(fileName)
		})

		ArchiveWithDomainRemover(archive, func(domain string) {
			fileName := filepath.Join(archive.name, domain)
			removeError := os.RemoveAll(fileName)
			if nil != removeError {
				NotifierSendError(archive.notifier, removeError)
			}
		})
	})
}
