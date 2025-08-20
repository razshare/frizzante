package program

import tea "github.com/charmbracelet/bubbletea"

func Run[T tea.Model](m T) (T, error) {
	r, err := tea.NewProgram(m, tea.WithFPS(120)).Run()
	if err != nil {
		return m, err
	}

	if t, ok := r.(T); ok {
		return t, nil
	}

	return m, nil
}
