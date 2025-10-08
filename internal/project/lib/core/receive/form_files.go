package receive

import (
	"mime/multipart"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormFiles reads all form files associated with the given key and returns them as a slice.
//
// If there are no files associated with the key, FormFiles stores an empty slice into value.
func FormFiles(client *clients.Client, key string, value *[]multipart.FileHeader) bool {
	if client.Request.MultipartForm == MultipartByReader {
		client.Config.ErrorLog.Println("http: multipart handled by MultipartReader", stack.Trace())
		*value = make([]multipart.FileHeader, 0)
		return false
	}

	if client.Request.MultipartForm == nil {
		if err := client.Request.ParseMultipartForm(MaxFormSize); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			*value = make([]multipart.FileHeader, 0)
			return false
		}
	}

	if client.Request.MultipartForm != nil && client.Request.MultipartForm.File != nil {
		if headers := client.Request.MultipartForm.File[key]; len(headers) > 0 {
			*value = make([]multipart.FileHeader, len(headers))
			for index, header := range headers {
				(*value)[index] = *header
			}
			return true
		}
	}

	*value = make([]multipart.FileHeader, 0)
	return false
}
