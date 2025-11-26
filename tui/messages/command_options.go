package messages

type CommandOptions struct {
	DisabledStdin bool
	DisableStdout bool
	DisableStderr bool
	DirectoryName string
	Environment   []string
	Program       string
	Args          []string
}
