package files

import (
	"bytes"
	"os"
)

func NewFileReader(name string) (reader *bytes.Reader, info os.FileInfo, err error) {
	var file *os.File
	if file, err = os.Open(name); err != nil {
		return
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()

	if info, err = file.Stat(); err != nil {
		return
	}

	buf := make([]byte, info.Size())
	if _, err = file.Read(buf); err != nil {
		return
	}

	if err = file.Close(); err != nil {
		return
	}

	return bytes.NewReader(buf), info, nil
}
