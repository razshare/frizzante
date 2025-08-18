package confirm

import (
	"errors"
	"fmt"
	"github.com/razshare/frizzante/stack"
	"github.com/razshare/frizzante/tui/program"
)

func Send(defaultValue bool, prompt string) (bool, error) {
	model, err := program.Run(&Model{
		Prompt:       prompt,
		DefaultValue: defaultValue,
		Confirmed:    defaultValue,
	})
	if err != nil {
		return false, errors.New(err.Error() + "\n" + stack.Trace())
	}
	return model.Confirmed, nil
}

func Sendf(defaultValue bool, format string, vars ...any) (bool, error) {
	return Send(defaultValue, fmt.Sprintf(format, vars...))
}
