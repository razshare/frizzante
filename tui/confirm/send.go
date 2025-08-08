package confirm

import "main/program"

func Send(prompt string, defaultValue bool) (bool, error) {
	model := Model{Prompt: prompt, DefaultValue: defaultValue, Confirmed: defaultValue}
	result, err := program.Run(model)
	if err != nil {
		return false, err
	}
	return result.Confirmed, nil
}
