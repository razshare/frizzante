package receive

import (
	"mime/multipart"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormFile reads the first form file associated with the given key and returns it.
func FormFile(client *client.Client, key string) MultipartFormFile {
	if !client.Parsed {
		Parse(client)
	}

	multipartForm := client.Request.MultipartForm
	if multipartForm != nil && multipartForm.File != nil {
		if headers := multipartForm.File[key]; len(headers) > 0 {
			var header *multipart.FileHeader
			if header = headers[0]; header == nil {
				client.Config.ErrorLog.Println("file header not found", stack.Trace())
				return MultipartFormFile{
					FileHeader: multipart.FileHeader{
						Header: map[string][]string{},
					},
				}
			}

			var err error
			var file multipart.File
			if file, err = header.Open(); err != nil {
				client.Config.ErrorLog.Println(err, stack.Trace())
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
