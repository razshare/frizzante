package receive

import (
	"mime/multipart"

	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

var MultipartByReader = &multipart.Form{
	Value: make(map[string][]string),
	File:  make(map[string][]*multipart.FileHeader),
}

// FormFile reads the first form file associated with the given key and returns it.
func FormFile(client *_client.Client, key string) MultipartFormFile {
	if client.Request.MultipartForm == MultipartByReader {
		client.Config.ErrorLog.Println("http: multipart handled by MultipartReader", stack.Trace())
		return MultipartFormFile{
			FileHeader: multipart.FileHeader{
				Header: map[string][]string{},
			},
		}
	}

	if client.Request.MultipartForm == nil {
		if err := client.Request.ParseMultipartForm(MaxFormSize); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return MultipartFormFile{
				FileHeader: multipart.FileHeader{
					Header: map[string][]string{},
				},
			}
		}
	}

	if client.Request.MultipartForm != nil && client.Request.MultipartForm.File != nil {
		if headers := client.Request.MultipartForm.File[key]; len(headers) > 0 {
			file, err := headers[0].Open()
			if err != nil {
				client.Config.ErrorLog.Println(err, stack.Trace())
				return MultipartFormFile{
					FileHeader: multipart.FileHeader{
						Header: map[string][]string{},
					},
				}
			}

			return MultipartFormFile{
				File: file,
				FileHeader: multipart.FileHeader{
					Header: map[string][]string{},
				},
			}
		}
	}

	return MultipartFormFile{
		FileHeader: multipart.FileHeader{
			Header: map[string][]string{},
		},
	}
}
