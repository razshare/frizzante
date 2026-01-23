package servers

// Stop stops the server.
func Stop(server *Server) (err error) {
	err = server.Close()
	return
}
