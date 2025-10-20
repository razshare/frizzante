package files

import "strings"

// FindWithPrefix finds all files with a given prefix within a directory recursively.
func FindWithPrefix(from string, prefix string) ([]string, error) {
	names, err := ReadDirectory(from)
	if err != nil {
		return nil, err
	}

	filtered := make([]string, 0)
	for _, name := range names {
		if strings.HasPrefix(name, prefix) {
			filtered = append(filtered, name)
		}
	}

	return filtered, nil
}
