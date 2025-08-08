package guard

import "github.com/razshare/frizzante/client"

type Guard struct {
	Name    string
	Handler func(c *client.Client, allow func())
	Tags    []string
}
