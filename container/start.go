package container

import (
	"strings"
	"sync"
)

// Start starts the container.
func Start(c *Container) {
	script := ProduceScript(
		c.Config.Efs,
		strings.ReplaceAll(c.Config.Root, "\\", "/"),
		strings.ReplaceAll(c.Config.Script, "\\", "/"),
		c.Config.InfoLog,
		c.Config.ErrorLog,
	)

	document := ProduceDocument(
		c.Config.Efs,
		strings.ReplaceAll(c.Config.Document, "\\", "/"),
		c.Config.InfoLog,
		c.Config.ErrorLog,
	)

	runtime := ProduceRuntime(
		c.Config.Parallels,
		c.Config.InfoLog,
	)

	program := ProduceProgram(
		c.Config.Script,
		script.Value,
		c.Config.Parallels,
		c.Config.InfoLog,
		c.Config.ErrorLog,
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
