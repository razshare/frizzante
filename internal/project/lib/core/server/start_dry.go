//go:build dry

package server

// Start starts a dry server from a configuration.
func Start(server *Server) {
	server.InfoLog.Println("server running dry")
}
