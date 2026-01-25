package embeds

import (
	"embed"
	"errors"
	"io"
	"io/fs"
)

// ReadFileInChunks reads a file in chunks of a set maximum size.
func ReadFileInChunks(efs embed.FS, from string, max int, call func([]byte) error) (err error) {
	var file fs.File
	if file, err = efs.Open(from); err != nil {
		return
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()
	buf := make([]byte, max)
	var size int
	for {
		size, err = file.Read(buf)
		if size > 0 {
			if err = call(buf[:size]); err != nil {
				return
			}
		}
		if errors.Is(err, io.EOF) {
			err = nil
			return
		}
	}
}
