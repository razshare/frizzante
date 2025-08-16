package files

import "os"

func Move(from string, to string) error {
	if IsDirectory(from) {
		err := CopyDirectory(from, to)
		if err != nil {
			return err
		}
	} else {
		err := CopyFile(from, to)
		if err != nil {
			return err
		}
	}

	err := os.RemoveAll(from)
	if err != nil {
		return err
	}

	return nil
}
