package receive

import (
	"mime/multipart"

	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// File returns the first file for the provided form key.
func (form *MultipartForm) File(key string) (file multipart.File, header *multipart.FileHeader, ok bool) {
	multipartForm := form.Client.Request.MultipartForm
	if multipartForm != nil && multipartForm.File != nil {
		if headers := multipartForm.File[key]; len(headers) > 0 {
			var err error
			header = headers[0]
			if file, err = header.Open(); err != nil {
				form.Client.Config.ErrorLog.Println(err, stack.Trace())
				return
			}
			ok = true
			return
		}
	}
	return
}
