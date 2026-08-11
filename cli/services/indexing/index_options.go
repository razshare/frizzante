package indexing

import "context"

type IndexOptions struct {
	Context     context.Context
	Address     string
	Depth       int
	OnProgress  func(current int, maximum int)
	StickToHost bool
}
