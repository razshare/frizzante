package npmselect

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/navigate"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func InitSearch() *search.Search {
	input := textinput.New()
	input.Width = 80
	
	return &search.Search{
		Active:   false,
		Choices:  []search.Choice{},
		Filtered: []search.Choice{},
		Input:    input,
	}
}

func InitViewport() *viewport.Viewport {
	return &viewport.Viewport{
		Visible: 6,
		Start:   0,
		Cursor:  0,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		if k.Type == tea.KeyCtrlC {
			return m, tea.Interrupt
		}
		
		if k.Type == tea.KeyEsc {
			m.Quitting = true
			m.Selected = []string{}
			return m, tea.Quit
		}

		if k.Type == tea.KeyEnter {
			if len(m.Selected) == 0 && len(m.Search.Filtered) > 0 {
				val := m.Search.Filtered[m.Viewport.Cursor].Id
				m.Selected = append(m.Selected, val)
			}
			m.Confirmed = true
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

		if k.Type == tea.KeyUp || k.Type == tea.KeyCtrlP {
			navigate.Apply(m.Search, m.Viewport, -1)
			return m, nil
		}

		if k.Type == tea.KeyDown || k.Type == tea.KeyCtrlN || k.Type == tea.KeyTab {
			navigate.Apply(m.Search, m.Viewport, 1)
			return m, nil
		}
		
		if k.String() == "ctrl+a" {
			if len(m.Search.Filtered) > 0 {
				if len(m.Selected) == len(m.Search.Filtered) {
					m.Selected = []string{}
				} else {
					m.Selected = []string{}
					for _, choice := range m.Search.Filtered {
						m.Selected = append(m.Selected, choice.Id)
					}
				}
			}
			return m, nil
		}

		// Handle search input
		if len(k.String()) == 1 || k.Type == tea.KeyBackspace || k.Type == tea.KeyCtrlH {
			if !m.Search.Active {
				m.Search.Active = true
				m.Search.Input.Focus()
			}
			
			var cmd tea.Cmd
			prevValue := m.Search.Input.Value()
			m.Search.Input, cmd = m.Search.Input.Update(k)
			
			if m.Search.Input.Value() != prevValue {
				m.LastQuery = m.Search.Input.Value()
				if m.DebounceTimer != nil {
					m.DebounceTimer.Stop()
				}
				// tea.Cmd returns a debounced search message after the specified delay
				return m, tea.Batch(cmd, func() tea.Msg {
					time.Sleep(500 * time.Millisecond)
					return DebouncedSearchMsg{Query: m.Search.Input.Value()}
				})
			}
			return m, cmd
		}

	case DebouncedSearchMsg:
		if k.Query != "" && k.Query == m.Search.Input.Value() {
			m.Loading = true
			return m, PerformSearch(k.Query)
		}

	case SearchResultMsg:
		m.Loading = false
		m.Error = k.Error
		m.Packages = k.Packages
		choices := make([]search.Choice, len(k.Packages))
		for i, pkg := range k.Packages {
			id := pkg.Name
			if pkg.Version != "" {
				id = fmt.Sprintf("%s@%s", pkg.Name, pkg.Version)
			}
			description := pkg.Description
			if len(description) > 50 {
				if 50 <= 3 {
					description = "..."
				} else {
					description = description[:50-3] + "..."
				}
			}
			choices[i] = search.Choice{
				Id:          id,
				Description: description,
			}
		}
		m.Search.Choices = choices
		m.Search.Filtered = choices
		m.Viewport.Cursor = 0
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder
	sb.Grow(1024)

	sb.WriteString(config.Styles.Menu.Render("Search NPM packages"))

	if m.Search.Input.Value() != "" {
		sb.WriteString(config.Styles.UserInput.Render(" ⁋/" + m.Search.Input.Value()))
	} else {
		sb.WriteString(config.Styles.UserGuide.Render(" ⁋/type to search"))
	}

	sb.WriteString("\n")

	footerText := "↑ up • ↓ down • space select • enter install • ctrl+a toggle all • esc cancel"

	// Handle status cases with early returns
	var statusContent string
	var isLoading, isError, isNoResults bool

	if m.Loading {
		statusContent = "⌛ searching..."
		isLoading = true
	} else if m.Error != nil {
		statusContent = fmt.Sprintf("✗ Error: %v", m.Error)
		isError = true
	} else if len(m.Search.Filtered) == 0 && m.Search.Input.Value() != "" {
		statusContent = "ⓘ  no packages found"
		isNoResults = true
	}

	if statusContent != "" {
		var statusSb strings.Builder
		statusSb.WriteString(config.Styles.Menu.Render("│"))
		
		if isLoading || isNoResults {
			statusSb.WriteString(config.Styles.UserGuide.PaddingLeft(1).Render(statusContent))
		} else if isError {
			statusSb.WriteString(config.Styles.Status(config.Colors.Error).PaddingLeft(1).Render(statusContent))
		}
		
		statusSb.WriteString("\n")
		statusSb.WriteString(config.Styles.UserGuide.Render(footerText))
		sb.WriteString(statusSb.String())
		return sb.String()
	}

	filtered := len(m.Search.Filtered)
	if filtered > 0 {
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
			choice := m.Search.Filtered[i]
			
			if m.Viewport.Cursor == i {
				if slices.Contains(m.Selected, choice.Id) {
					sb.WriteString(config.Styles.Selected.Render("● " + choice.Id))
				} else {
					sb.WriteString(config.Styles.Selected.Render("◉ " + choice.Id))
				}
				
				if choice.Description != "" {
					sb.WriteString(config.Styles.UserGuide.Render("  ⇢  " + choice.Description))
				}
			} else if slices.Contains(m.Selected, choice.Id) {
				sb.WriteString(config.Styles.Item.Render("● " + choice.Id))
			} else {
				sb.WriteString(config.Styles.Item.Render("○ " + choice.Id))
			}
			sb.WriteString("\n")
		}
		
		if height < filtered {
			sb.WriteString(config.Styles.Menu.Render("│"))
			sb.WriteString(config.Styles.Status(config.Colors.Muted).Render("↓ more below"))
			sb.WriteString("\n")
		}
	}

	sb.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter install"))
	if len(m.Selected) > 0 {
		sb.WriteString(config.Styles.UserGuide.Render(fmt.Sprintf(" (%d selected)", len(m.Selected))))
	}
	sb.WriteString(config.Styles.UserGuide.Render(" • ctrl+a toggle all • esc cancel"))

	return sb.String()
}

func PerformSearch(query string) tea.Cmd {
	return func() tea.Msg {
		if query == "" {
			return SearchResultMsg{Packages: []PackageInfo{}}
		}

		encodedQuery := url.QueryEscape(query)
		apiUrl := fmt.Sprintf("https://registry.npmjs.org/-/v1/search?text=%s&size=20", encodedQuery)

		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			return SearchResultMsg{Error: err}
		}

		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return SearchResultMsg{Error: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return SearchResultMsg{Error: fmt.Errorf("npm registry returned status %d", resp.StatusCode)}
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return SearchResultMsg{Error: err}
		}

		var searchResult SearchResponse
		err = json.Unmarshal(body, &searchResult)
		if err != nil {
			return SearchResultMsg{Error: err}
		}

		packages := make([]PackageInfo, 0, len(searchResult.Objects))
		for _, obj := range searchResult.Objects {
			packages = append(packages, obj.Package)
		}

		return SearchResultMsg{Packages: packages}
	}
}

func ExtractPackageNames(selected []string) []string {
	names := make([]string, 0, len(selected))
	for _, id := range selected {
		// Remove version suffix if present (e.g., "package@1.0.0" -> "package")
		if idx := strings.IndexByte(id, '@'); idx != -1 {
			names = append(names, id[:idx])
		} else {
			names = append(names, id)
		}
	}
	return names
}
