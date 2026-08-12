package messages

import "strings"

type CommandOptions struct {
	DisabledStdin bool
	DisableStdout bool
	DisableStderr bool
	DirectoryName string
	Environment   []string
	Program       string
	Args          []string
	Channels      CommandChannels
	StdoutBuilder *strings.Builder
	StderrBuilder *strings.Builder
}
