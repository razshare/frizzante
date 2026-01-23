package servers

type Channels struct {
	Started chan struct{}
	Ended   chan struct{}
}
