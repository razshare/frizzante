package container

// New starts a new app.
func New(conf Config) *Container {
	return &Container{
		Config:   conf,
		Channels: Channels{},
	}
}
