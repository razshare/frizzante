package singleselect

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func (model *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (model *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch messageLocal := msg.(type) {
	case tea.KeyMsg:
		key := messageLocal.String()
		switch key {
		case "ctrl+c":
			return model, tea.Quit
		case "esc":
			if model.Searching {
				model.ResetSearch()
			}
			return model, nil
		case "enter":
			if len(model.FilteredChoices) > 0 {
				model.Selected = model.FilteredChoices[model.Cursor]
				return model, tea.Quit
			}
		case "up", "ctrl+p":
			model.Navigate(-1)
			return model, nil
		case "down", "ctrl+n", "tab":
			model.Navigate(1)
			return model, nil
		case "backspace", "ctrl+h":
			if model.Searching {
				return model.HandleSearchInput(messageLocal)
			}
		default:
			if len(key) == 1 || key == "space" {
				if !model.Searching {
					model.Searching = true
					model.SearchInput.Focus()
				}
				return model.HandleSearchInput(messageLocal)
			}
		}
	}
	return model, nil
}

func (model *Model) View() string {
	var s strings.Builder
	s.Grow(1024)

	s.WriteString(config.Styles.Title.Render(model.Prompt))
	s.WriteString(" ")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(config.Colors.Secondary)).Render("[type to search]"))
	s.WriteString(": ")
	s.WriteString(model.SearchInput.Value())
	s.WriteString("\n")

	choiceCount := len(model.FilteredChoices)
	if choiceCount == 0 {
		s.WriteString("  No matches found\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(config.Colors.Muted)).Render("(esc to clear search, ctrl + c to quit)"))
		return s.String()
	}

	viewportEnd := model.ViewportStart + model.MaxVisible
	if viewportEnd > choiceCount {
		viewportEnd = choiceCount
	}

	if model.ViewportStart > 0 {
		s.WriteString(config.Styles.Status(config.Colors.Muted).Render("    ↑ more above"))
		s.WriteString("\n")
	}

	for i := model.ViewportStart; i < viewportEnd; i++ {
		if i == model.Cursor {
			s.WriteString(config.Styles.Selected.Render("▶ " + model.FilteredChoices[i]))
		} else {
			s.WriteString(config.Styles.Item.Render("  " + model.FilteredChoices[i]))
		}
		s.WriteString("\n")
	}

	if viewportEnd < choiceCount {
		s.WriteString(config.Styles.Status(config.Colors.Muted).Render("    ↓ more below"))
		s.WriteString("\n")
	}

	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(config.Colors.Muted)).Render("(↑/↓ Navigate, enter to select, esc to clear search, ctrl + c to quit)"))
	return s.String()
}

func (model *Model) ResetSearch() {
	model.SearchInput.SetValue("")
	model.Searching = false
	model.FilteredChoices = model.Choices
	model.Cursor = 0
	model.ViewportStart = 0
}

func (model *Model) HandleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	prevValue := model.SearchInput.Value()
	var cmd tea.Cmd
	model.SearchInput, cmd = model.SearchInput.Update(msg)

	if newValue := model.SearchInput.Value(); newValue != prevValue {
		if newValue == "" {
			model.ResetSearch()
		} else {
			model.FilterChoices()
		}
	}
	return model, cmd
}

func (model *Model) Navigate(direction int) {
	count := len(model.FilteredChoices)
	if count == 0 {
		return
	}

	model.Cursor = (model.Cursor + direction + count) % count

	if model.Cursor < model.ViewportStart {
		model.ViewportStart = model.Cursor
	} else if model.Cursor >= model.ViewportStart+model.MaxVisible {
		model.ViewportStart = model.Cursor - model.MaxVisible + 1
	}
}

func (model *Model) FilterChoices() {
	searchTerm := strings.ToLower(model.SearchInput.Value())
	if searchTerm == "" {
		model.FilteredChoices = model.Choices
	} else {
		filtered := make([]string, 0, len(model.Choices)/2)
		for _, choice := range model.Choices {
			if strings.Contains(strings.ToLower(choice), searchTerm) {
				filtered = append(filtered, choice)
			}
		}
		model.FilteredChoices = filtered
	}
	model.Cursor = 0
	model.ViewportStart = 0
}
