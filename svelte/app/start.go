package app

import (
	"github.com/razshare/frizzante/server"
	"strings"
	"sync"
)

// Start starts the app and its server.
func Start(c *Config) {
	script := ProduceScript(
		c.Server.Efs,
		strings.ReplaceAll(c.Root, "\\", "/"),
		strings.ReplaceAll(c.Script, "\\", "/"),
		c.Server.InfoLog,
		c.Server.ErrorLog,
	)

	document := ProduceDocument(
		c.Server.Efs,
		strings.ReplaceAll(c.Document, "\\", "/"),
		c.Server.InfoLog,
		c.Server.ErrorLog,
	)

	runtime := ProduceRuntime(
		c.Parallels,
		c.Server.InfoLog,
	)

	program := ProduceProgram(
		c.Script,
		script.Value,
		c.Parallels,
		c.Server.InfoLog,
		c.Server.ErrorLog,
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

	defer func() { stop <- 0 }()
	server.Start(c.Server)
}
