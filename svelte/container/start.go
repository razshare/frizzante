package container

import (
	"strings"
	"sync"
)

// Start starts ssr producers and the server.
func Start(c *Container) {
	script := ProduceScript(
		c.Efs,
		strings.ReplaceAll(c.Root, "\\", "/"),
		strings.ReplaceAll(c.Script, "\\", "/"),
		c.ErrorLog,
	)

	document := ProduceDocument(
		c.Efs,
		strings.ReplaceAll(c.Document, "\\", "/"),
		c.ErrorLog,
	)

	runtime := ProduceRuntime(
		c.Parallels,
	)

	program := ProduceProgram(
		c.Script,
		script.Value,
		c.Parallels,
		c.ErrorLog,
	)

	stop := make(chan any, 1)

	go func() {
		<-stop
		var group sync.WaitGroup
		group.Add(4)
		go func() { group.Done(); document.Stop <- 0 }()
		go func() { group.Done(); script.Stop <- 0 }()
		go func() { group.Done(); program.Stop <- 0 }()
		go func() { group.Done(); runtime.Stop <- 0 }()
		group.Wait()
	}()

	c.Channels.Script = script.Value
	c.Channels.Document = document.Value
	c.Channels.Program = program.Value
	c.Channels.Runtime = runtime.Value
	c.Channels.Stop = stop
}
