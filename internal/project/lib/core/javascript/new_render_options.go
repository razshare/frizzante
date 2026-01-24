package javascript

import "log"

type NewRenderOptions struct {
	Data     []byte
	Server   string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
