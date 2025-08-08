package server

import (
	"github.com/razshare/frizzante/container"
	"github.com/razshare/frizzante/guard"
	"github.com/razshare/frizzante/route"
	"log"
	"net/http"
)

type Channels struct {
	Stop chan any
}

type Config struct {
	Container   container.Config
	Guards      []guard.Guard
	Routes      []route.Route
	Http        *http.Server
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
	Channels    Channels
	SecureAddr  string
	Certificate string
	Key         string
}
