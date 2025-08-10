package multiselect

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/navigate"
	"github.com/razshare/frizzante/tui/search"
	"os"
	"strings"
)

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		if k.Type == tea.KeyCtrlC {
			os.Exit(0)
		}

		if k.Type == tea.KeyEnter {
			return model, tea.Quit
		}

		if k.Type == tea.KeySpace {
			if len(model.Search.Filtered) > 0 {
				if model.Selected[model.Viewport.Cursor] {
					delete(model.Selected, model.Viewport.Cursor)
				} else {
					model.Selected[model.Viewport.Cursor] = true
				}
			}
			return model, nil
		}

		if k.Type == tea.KeyEsc {
			if model.Search.Active {
				search.Reset(model.Search, model.Viewport)
			}
			return model, nil
		}

		if k.Type == tea.KeyUp || k.Type == tea.KeyCtrlP {
			navigate.Apply(model.Search, model.Viewport, -1)
			return model, nil
		}

		if k.Type == tea.KeyDown || k.Type == tea.KeyCtrlN || k.Type == tea.KeyTab {
			navigate.Apply(model.Search, model.Viewport, 1)
			return model, nil
		}

		if k.Type == tea.KeyBackspace || k.Type == tea.KeyCtrlH {
			if model.Search.Active {
				return model, search.Apply(model.Search, model.Viewport, k)
			}
		}

		if len(k.String()) == 1 {
			if !model.Search.Active {
				model.Search.Active = true
				model.Search.Input.Focus()
			}
			return model, search.Apply(model.Search, model.Viewport, k)
		}
	}
	return model, nil
}

func (model Model) View() string {
	var sbguide strings.Builder
	sbguide.WriteString("(")
	sbguide.WriteString("↑/↓ = navigate, [Enter] = select")
	if model.Search.Active {
		sbguide.WriteString(", [Esc] = clear")
	}
	sbguide.WriteString(")")

	var sb strings.Builder
	sb.Grow(1024)

	sb.WriteString(config.Styles.Title.Render(model.Prompt))
	sb.WriteString(" ")
	sb.WriteString(config.Styles.Suggestion.Render("[type to search]"))
	sb.WriteString(": ")
	sb.WriteString(config.Styles.UserInput.Render(model.Search.Input.Value()))
	sb.WriteString("\n")

	choices := len(model.Search.Filtered)
	if choices == 0 {
		sb.WriteString("  No matches found\n")
		sb.WriteString(config.Styles.UserGuide.Render(sbguide.String()))
		return sb.String()
	}

	height := model.Viewport.Start + model.Viewport.Visible
	if height > choices {
		height = choices
	}

	if model.Viewport.Start > 0 {
		sb.WriteString(config.Styles.Status(config.Colors.Muted).Render("    ↑ more above"))
		sb.WriteString("\n")
	}

	for i := model.Viewport.Start; i < height; i++ {
		var sbloc strings.Builder

		sbloc.WriteString(" ")
		if model.Viewport.Cursor == i {
			sbloc.WriteString(">")
		}

		sbloc.WriteString("[")

		if model.Selected[i] {
			sbloc.WriteString("✓")
		}

		sbloc.WriteString("]")

		sbloc.WriteString(model.Search.Choices[i])

		if model.Viewport.Cursor == i {
			sb.WriteString(config.Styles.Selected.Render(sbloc.String()) + "\n")
		} else {
			sb.WriteString(config.Styles.Item.Render(sbloc.String()) + "\n")
		}
	}

	if height < choices {
		sb.WriteString(config.Styles.Status(config.Colors.Muted).Render("    ↓ more below"))
		sb.WriteString("\n")
	}

	sb.WriteString(config.Styles.UserGuide.Render(sbguide.String()))

	return sb.String()
}
