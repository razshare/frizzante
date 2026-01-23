package javascript

import "log"

type NewRenderFunctionOptions struct {
	Data     []byte
	Server   string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
