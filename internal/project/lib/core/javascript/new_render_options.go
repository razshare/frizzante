package javascript

import "log"

type NewRenderOptions struct {
	FindSource func() (source string, err error)
	Server     string
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
}
