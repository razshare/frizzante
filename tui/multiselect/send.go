package multiselect

import (
	"main/program"
)

func Send(prompt string, options []string) ([]string, error) {
	m := Model{Choices: options, Selected: make(map[int]bool), Prompt: prompt}
	result, err := program.Run(m)
	if err != nil {
		return nil, err
	}
	var selections []string
	for i, choice := range result.Choices {
		if result.Selected[i] {
			selections = append(selections, choice)
		}
	}
	return selections, nil
}
