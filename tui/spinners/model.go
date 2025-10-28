package spinners

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
)

func (model *Model) Init() tea.Cmd {
	return model.Spinner.Tick
}

func (model *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch assert := message.(type) {
	case tea.KeyMsg:
		if assert.Type == tea.KeyCtrlC {
			return model, tea.Interrupt
		}
	}

	var cmd tea.Cmd
	model.Spinner, cmd = model.Spinner.Update(message)
	return model, cmd
}

func (model *Model) View() string {
	return fmt.Sprintf("\r%s %s", model.Spinner.View(), config.Styles.Menu.Render(model.Message))
}
