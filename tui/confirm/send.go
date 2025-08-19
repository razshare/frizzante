package confirm

import (
	"fmt"
	"github.com/razshare/frizzante/tui/program"
)

func Send(def bool, msg string) (bool, error) {
	model, err := program.Run(&Model{
		Prompt:       msg,
		DefaultValue: def,
		Confirmed:    def,
	})

	if err != nil {
		return false, err
	}

	return model.Confirmed, nil
}

func Sendf(def bool, format string, vars ...any) (bool, error) {
	return Send(def, fmt.Sprintf(format, vars...))
}
