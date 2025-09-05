package files

import "os"

func Move(from string, to string) (err error) {
	if IsDirectory(from) {
		if err = CopyDirectory(from, to); err != nil {
			return
		}
	} else {
		if err = CopyFile(from, to); err != nil {
			return
		}
	}

	if err = os.RemoveAll(from); err != nil {
		return
	}

	return nil
}
