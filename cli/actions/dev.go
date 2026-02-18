package actions

import (
	"os"

	"github.com/razshare/frizzante/cli/generations"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Dev(options DevOptions) (err error) {
	if !files.IsFile(".air.toml") {
		if err = generations.AirConfig(generations.AirConfigOptions{
			Efs: options.Efs,
		}); err != nil {
			return
		}
	}
	if !messages.Command(messages.CommandOptions{
		Environment: append(os.Environ(), "DEV=1"),
		Program:     options.Air,
	}) {
		messages.Error("air failed")
	}
	return
}
