package program

import tea "github.com/charmbracelet/bubbletea"

func Run[T tea.Model](model T) (T, error) {
	prog, err := tea.NewProgram(model, tea.WithFPS(120)).Run()
	if err != nil {
		return model, err
	}

	if cast, ok := prog.(T); ok {
		return cast, nil
	}

	return model, nil
}
