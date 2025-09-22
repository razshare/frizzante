package receive

import (
	"mime/multipart"

	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// File returns the first file for the provided form key.
func (form *MultipartForm) File(key string) MultipartFormFile {
	multipartForm := form.Client.Request.MultipartForm
	if multipartForm != nil && multipartForm.File != nil {
		if headers := multipartForm.File[key]; len(headers) > 0 {
			var header *multipart.FileHeader
			if header = headers[0]; header == nil {
				form.Client.Config.ErrorLog.Println("file header not found", stack.Trace())
				return MultipartFormFile{
					FileHeader: multipart.FileHeader{
						Header: map[string][]string{},
					},
				}
			}

			var err error
			var file multipart.File
			if file, err = header.Open(); err != nil {
				form.Client.Config.ErrorLog.Println(err, stack.Trace())
				return MultipartFormFile{
					FileHeader: *header,
				}
			}
			return MultipartFormFile{
				File:       file,
				FileHeader: *header,
			}
		}
	}

	return MultipartFormFile{
		FileHeader: multipart.FileHeader{
			Header: map[string][]string{},
		},
	}
}
