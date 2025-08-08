package multiselect

import (
	tea "github.com/charmbracelet/bubbletea"
	"main/config"
)

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch messageLocal := msg.(type) {
	case tea.KeyMsg:
		switch messageLocal.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		case "up", "k":
			if model.Cursor > 0 {
				model.Cursor--
			}
		case "down", "j":
			if model.Cursor < len(model.Choices)-1 {
				model.Cursor++
			}
		case " ":
			if model.Selected[model.Cursor] {
				delete(model.Selected, model.Cursor)
			} else {
				model.Selected[model.Cursor] = true
			}
		case "enter":
			return model, tea.Quit
		}
	}
	return model, nil
}

func (model Model) View() string {
	s := config.Styles.Title.Render(model.Prompt) + "\n"
	s += config.Styles.Status(config.Colors.Info).Render("Use arrow keys to Navigate, space to select, enter to confirm") + "\n"
	for i, choice := range model.Choices {
		cursor := " "
		if model.Cursor == i {
			cursor = ">"
		}
		checked := " "
		if model.Selected[i] {
			checked = "✓"
		}
		line := cursor + " [" + checked + "] " + choice
		if model.Cursor == i {
			s += config.Styles.Selected.Render(line) + "\n"
		} else {
			s += config.Styles.Item.Render(line) + "\n"
		}
	}
	return s
}
