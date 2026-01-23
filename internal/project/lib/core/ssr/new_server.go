package ssr

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/servers"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// NewServer creates a new server.
func NewServer() (server *servers.Server) {
	errorLog := log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime)
	limitString := os.Getenv("FRIZZANTE_JS_RUNTIME_LIMIT")
	limit := 1
	if limitString != "" {
		var err error
		if limit, err = strconv.Atoi(limitString); err != nil {
			limit = 1
			errorLog.Printf(
				"could not parse FRIZZANTE_JS_RUNTIME_LIMIT environment variable; reverting back to limit 1; received %s %s",
				limitString,
				stack.Trace(),
			)
		}
	}
	return &servers.Server{
		InfoLog:        log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime),
		SecureAddr:     "0.0.0.0:8383",
		Cors:           http.NewCrossOriginProtection(),
		Addr:           "0.0.0.0:8080",
		Handler:        http.NewServeMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 2097152, // 2MB
		ErrorLog:       errorLog,
		Render:         NewFunction(int64(limit)),
		Channels: servers.Channels{
			Started: make(chan struct{}, 1),
			Ended:   make(chan struct{}, 1),
		},
	}
}
