package files

import "strings"

// FindWithSuffix finds all files with a given suffix within a directory recursively.
func FindWithSuffix(from string, suffix string) ([]string, error) {
	names, err := ReadDirectory(from)
	if err != nil {
		return nil, err
	}

	filtered := make([]string, 0)
	for _, name := range names {
		if strings.HasSuffix(name, suffix) {
			filtered = append(filtered, name)
		}
	}

	return filtered, nil
}
