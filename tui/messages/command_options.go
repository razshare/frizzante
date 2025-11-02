package messages

type CommandOptions struct {
	DisabledStdin bool
	DisableStdout bool
	DisableStderr bool
	Dir           string
	Env           []string
	Name          string
	Args          []string
}
