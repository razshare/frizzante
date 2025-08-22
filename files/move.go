package files

import "os"

func Move(from string, to string) error {
	if IsDirectory(from) {
		if err := CopyDirectory(from, to); err != nil {
			return err
		}
	} else {
		if err := CopyFile(from, to); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(from); err != nil {
		return err
	}

	return nil
}
