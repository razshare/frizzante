package singleselect

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/navigate"
	"github.com/razshare/frizzante/tui/search"
	"os"
	"slices"
	"strings"
)

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		if k.Type == tea.KeyCtrlC {
			os.Exit(0)
		}

		if k.Type == tea.KeyEnter {
			if len(m.Search.Filtered) > 0 {
				m.Selected = m.Search.Filtered[m.Viewport.Cursor]
				return m, tea.Quit
			}
		}

		if k.Type == tea.KeyEsc {
			if m.Search.Active {
				search.Reset(m.Search, m.Viewport)
			}
			return m, nil
		}

		if k.Type == tea.KeyUp || k.Type == tea.KeyCtrlP {
			navigate.Apply(m.Search, m.Viewport, -1)
			return m, nil
		}

		if k.Type == tea.KeyDown || k.Type == tea.KeyCtrlN || k.Type == tea.KeyTab {
			navigate.Apply(m.Search, m.Viewport, 1)
			return m, nil
		}

		if k.Type == tea.KeyBackspace || k.Type == tea.KeyCtrlH {
			if m.Search.Active {
				return m, search.Apply(m.Search, m.Viewport, k)
			}
		}

		if len(k.String()) == 1 {
			if !m.Search.Active {
				m.Search.Active = true
				m.Search.Input.Focus()
			}
			return m, search.Apply(m.Search, m.Viewport, k)
		}
	}

	return m, nil
}

func (m *Model) View() string {
	var sbguide strings.Builder
	sbguide.WriteString("(")
	sbguide.WriteString("↑/↓ = navigate, [Enter] = select")
	if m.Search.Active {
		sbguide.WriteString(", [Esc] = clear")
	}
	sbguide.WriteString(")")

	var sb strings.Builder
	sb.Grow(1024)

	sb.WriteString(config.Styles.Title.Render(m.Prompt))
	sb.WriteString(" ")
	sb.WriteString(config.Styles.Suggestion.Render("[type to search]"))
	sb.WriteString(": ")
	sb.WriteString(config.Styles.UserInput.Render(m.Search.Input.Value()))
	sb.WriteString("\n")

	filtered := len(m.Search.Filtered)
	if filtered == 0 {
		sb.WriteString("  No matches found\n")
		sb.WriteString(config.Styles.UserGuide.Render(sbguide.String()))
		return sb.String()
	}

	height := m.Viewport.Start + m.Viewport.Visible
	if height > filtered {
		height = filtered
	}

	if m.Viewport.Start > 0 {
		sb.WriteString(config.Styles.Status(config.Colors.Muted).Render("↑ more above"))
		sb.WriteString("\n")
	}

	for i := m.Viewport.Start; i < height; i++ {
		if i == m.Viewport.Cursor {
			sb.WriteString(config.Styles.Selected.Render("● " + m.Search.Filtered[i]))
			j := slices.Index(m.Search.Choices, m.Search.Filtered[i])
			if j >= 0 && m.Search.Descriptions[j] != "" {
				sb.WriteString(config.Styles.UserGuide.Render("  ⇢  " + m.Search.Descriptions[j]))
			}
		} else {
			sb.WriteString(config.Styles.Item.Render("○ " + m.Search.Filtered[i]))
		}
		sb.WriteString("\n")
	}

	if height < filtered {
		sb.WriteString(config.Styles.Status(config.Colors.Muted).Render("↓ more below"))
		sb.WriteString("\n")
	}

	sb.WriteString(config.Styles.UserGuide.Render(sbguide.String()))

	return sb.String()
}
