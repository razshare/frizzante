package receive

import (
	"mime/multipart"
	"net/url"

	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
)

type MultipartFormFile struct {
	multipart.File
	multipart.FileHeader
}

type MultipartForm struct {
	url.Values
	Client *_client.Client
}
