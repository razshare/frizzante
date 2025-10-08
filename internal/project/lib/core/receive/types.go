package receive

import (
	"mime/multipart"
	"net/url"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

type MultipartFormFile struct {
	multipart.File
	multipart.FileHeader
}

type MultipartForm struct {
	url.Values
	Client *clients.Client
}
