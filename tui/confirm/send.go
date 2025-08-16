package confirm

import (
	"fmt"
	"github.com/razshare/frizzante/tui/program"
)

func Send(defaultValue bool, prompt string) (bool, error) {
	model := Model{Prompt: prompt, DefaultValue: defaultValue, Confirmed: defaultValue}
	result, err := program.Run(model)
	if err != nil {
		return false, err
	}
	return result.Confirmed, nil
}

func Sendf(defaultValue bool, format string, vars ...any) (bool, error) {
	return Send(defaultValue, fmt.Sprintf(format, vars...))
}
