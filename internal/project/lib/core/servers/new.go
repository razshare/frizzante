package servers

import (
	"log"
	"net/http"
	"time"
)

func New(errorLog *log.Logger, infoLog *log.Logger) *Server {
	return &Server{
		InfoLog:    infoLog,
		SecureAddr: "0.0.0.0:8383",
		Cors:       http.NewCrossOriginProtection(),
		Server: http.Server{
			Addr:           "0.0.0.0:8080",
			Handler:        http.NewServeMux(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 2097152, // 2MB,
			ErrorLog:       errorLog,
		},
	}
}
