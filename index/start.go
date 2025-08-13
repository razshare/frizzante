package index

import (
	"strings"
	"sync"
)

// Start starts csr producers and the server.
func Start(c *Config) {
	document := ProduceDocument(
		c.Efs,
		strings.ReplaceAll(c.Document, "\\", "/"),
		c.ErrorLog,
	)

	stop := make(chan any, 1)

	go func() {
		<-stop
		var group sync.WaitGroup
		group.Add(4)
		go func() { group.Done(); document.Stop <- 0 }()
		group.Wait()
	}()

	c.Channels.Document = document.Value
	c.Channels.Stop = stop
}
