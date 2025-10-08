package receive

import "mime/multipart"

type MultipartFormFile struct {
	multipart.File
	multipart.FileHeader
}
