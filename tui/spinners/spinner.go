package spinners

import tea "github.com/charmbracelet/bubbletea"

type Spinner struct {
	Model   *Model
	Program *tea.Program
	Done    chan struct{}
}
