package servers

type Channels struct {
	Start chan struct{}
	End   chan struct{}
}
