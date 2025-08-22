package embeds

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
)

func NewFileReader(efs embed.FS, n string) (reader *bytes.Reader, info os.FileInfo, err error) {
	var file fs.File
	if file, err = efs.Open(n); err != nil {
		return
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()

	info, _ = file.Stat()

	buf := make([]byte, info.Size())

	if _, err = file.Read(buf); err != nil {
		return
	}

	if err = file.Close(); err != nil {
		return
	}

	reader = bytes.NewReader(buf)

	return

}
