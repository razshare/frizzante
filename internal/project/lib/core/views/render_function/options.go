package render_function

import "log"

type Options struct {
	Data     []byte
	Server   string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
