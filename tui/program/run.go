package program

import tea "github.com/charmbracelet/bubbletea"

func Run[T tea.Model](model T) (T, error) {
	result, err := tea.NewProgram(model, tea.WithFPS(120)).Run()
	if err != nil {
		return model, err
	}

	if typed, ok := result.(T); ok {
		return typed, nil
	}

	return model, nil
}
