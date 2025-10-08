package receive

import (
	"mime/multipart"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

var MultipartByReader = &multipart.Form{
	Value: make(map[string][]string),
	File:  make(map[string][]*multipart.FileHeader),
}

// FormFile reads the first form file associated with the given key and returns it.
func FormFile(client *clients.Client, key string) MultipartFormFile {
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

// FormFileSlice reads all form files associated with the given key and returns them as a slice.
//
// If there are no files associated with the key, FormFileSlice stores an empty slice into value.
// If any file fails to open, FormFileSlice returns false.
func FormFileSlice(client *clients.Client, key string, value *[]MultipartFormFile) bool {
	if client.Request.MultipartForm == MultipartByReader {
		client.Config.ErrorLog.Println("http: multipart handled by MultipartReader", stack.Trace())
		return false
	}

	if client.Request.MultipartForm == nil {
		if err := client.Request.ParseMultipartForm(MaxFormSize); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}
	}

	*value = make([]MultipartFormFile, 0)

	if client.Request.MultipartForm != nil && client.Request.MultipartForm.File != nil {
		if headers := client.Request.MultipartForm.File[key]; len(headers) > 0 {
			for _, header := range headers {
				file, err := header.Open()
				if err != nil {
					client.Config.ErrorLog.Println(err, stack.Trace())
					return false
				}

				*value = append(*value, MultipartFormFile{
					File: file,
					FileHeader: multipart.FileHeader{
						Filename: header.Filename,
						Header:   header.Header,
						Size:     header.Size,
					},
				})
			}
		}
	}

	return true
}
