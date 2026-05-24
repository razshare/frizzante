package confirm

import (
	"github.com/razshare/frizzante/tui/program"
)

func Send(defaultValue bool, message string) (yes bool, err error) {
	var model *Model
	if model, err = program.Run(&Model{Prompt: message, Confirmed: defaultValue}); err != nil {
		return
	}
	yes = model.Confirmed
	return
}
