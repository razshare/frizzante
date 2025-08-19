package multiselect

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/navigate"
	"github.com/razshare/frizzante/tui/search"
	"slices"
	"strings"
)

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		if k.Type == tea.KeyCtrlC {
			return m, tea.Interrupt
		}

		if k.Type == tea.KeyEnter {
			if len(m.Selected) == 0 && len(m.Search.Filtered) > 0 {
				val := m.Search.Filtered[m.Viewport.Cursor].Id
				m.Selected = append(m.Selected, val)
			}
			return m, tea.Quit
		}

		if k.Type == tea.KeySpace {
			if len(m.Search.Filtered) > 0 {
				val := m.Search.Filtered[m.Viewport.Cursor].Id
				if slices.Contains(m.Selected, val) {
					if i := slices.Index(m.Selected, val); i >= 0 {
						m.Selected = append(m.Selected[:i], m.Selected[i+1:]...)
					}
				} else {
					m.Selected = append(m.Selected, val)
				}
			}
			return m, nil
		}

		if k.Type == tea.KeyEsc {
			if m.Search.Active {
				search.Reset(m.Search, m.Viewport)
			}

			m.Selected = make([]string, 0)
			return m, tea.Quit
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
	var sb strings.Builder
	sb.Grow(1024)

	sb.WriteString(config.Styles.Menu.Render(m.Prompt))

	if m.Search.Input.Value() != "" {
		sb.WriteString(config.Styles.UserInput.Render(" ⁋/" + m.Search.Input.Value()))
	} else {
		sb.WriteString(config.Styles.UserGuide.Render(" ⁋/type to search"))
	}

	sb.WriteString("\n")

	filtered := len(m.Search.Filtered)
	if filtered == 0 {
		sb.WriteString(config.Styles.Menu.Render("│"))
		sb.WriteString(config.Styles.UserGuide.PaddingLeft(1).Render("ⓘ  no matches found"))

		sb.WriteString("\n")

		sb.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter continue"))

		if m.Search.Active {
			sb.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
		} else {
			sb.WriteString(config.Styles.UserGuide.Render(" • esc back"))
		}

		return sb.String()
	}

	height := m.Viewport.Start + m.Viewport.Visible
	if height > filtered {
		height = filtered
	}

	if m.Viewport.Start > 0 {
		sb.WriteString(config.Styles.Menu.Render("│"))
		sb.WriteString(config.Styles.Status(config.Colors.Muted).Render("↑ more above"))
		sb.WriteString("\n")
	}

	for i := m.Viewport.Start; i < height; i++ {
		sb.WriteString(config.Styles.Menu.Render("│"))
		if m.Viewport.Cursor == i {
			if slices.Contains(m.Selected, m.Search.Filtered[i].Id) {
				sb.WriteString(config.Styles.Selected.Render("● " + m.Search.Filtered[i].Id))
			} else {
				sb.WriteString(config.Styles.Selected.Render("◉ " + m.Search.Filtered[i].Id))
			}

			j := slices.Index(m.Search.Choices, m.Search.Filtered[i])
			if j >= 0 && m.Search.Choices[j].Description != "" {
				sb.WriteString(config.Styles.UserGuide.Render("  ⇢  " + m.Search.Choices[j].Description))
			}
		} else if slices.Contains(m.Selected, m.Search.Filtered[i].Id) {
			sb.WriteString(config.Styles.Item.Render("● " + m.Search.Filtered[i].Id))
		} else {
			sb.WriteString(config.Styles.Item.Render("○ " + m.Search.Filtered[i].Id))
		}
		sb.WriteString("\n")
	}

	if height < filtered {
		sb.WriteString(config.Styles.Menu.Render("│"))
		sb.WriteString(config.Styles.Status(config.Colors.Muted).Render("↓ more below"))
		sb.WriteString("\n")
	}

	sb.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter continue"))

	if m.Search.Active {
		sb.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
	} else {
		sb.WriteString(config.Styles.UserGuide.Render(" • esc back"))
	}

	return sb.String()
}
