package confirm

import (
	"fmt"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/program"
)

func Send(defaultValue bool, prompt string) bool {
	model := Model{Prompt: prompt, DefaultValue: defaultValue, Confirmed: defaultValue}
	result, err := program.Run(model)
	if err != nil {
		messages.Fatal(err)
	}
	return result.Confirmed
}

func Sendf(defaultValue bool, format string, vars ...any) bool {
	return Send(defaultValue, fmt.Sprintf(format, vars...))
}
