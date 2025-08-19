package action

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/navigate"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/spinner"
	"github.com/razshare/frizzante/tui/viewport"
)

func InitNpmSearch() NpmSearchModel {
	input := textinput.New()
	input.Width = 80
	
	return NpmSearchModel{
		Search: &search.Search{
			Active:   false,
			Choices:  []search.Choice{},
			Filtered: []search.Choice{},
			Input:    input,
		},
		Viewport: &viewport.Viewport{
			Visible: 6,
			Start:   0,
			Cursor:  0,
		},
		Packages:      []NpmPackageInfo{},
		Selected:      []string{},
		Loading:       false,
		DebounceTimer: nil,
	}
}

func ConvertPackagesToChoices(packages []NpmPackageInfo) []search.Choice {
	choices := make([]search.Choice, len(packages))
	for i, pkg := range packages {
		id := pkg.Name
		if pkg.Version != "" {
			id = fmt.Sprintf("%s@%s", pkg.Name, pkg.Version)
		}
		choices[i] = search.Choice{
			Id:          id,
			Description: TruncateText(pkg.Description, 50),
		}
	}
	return choices
}

func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	if maxLength <= 3 {
		return "..."
	}
	return text[:maxLength-3] + "..."
}

func (m NpmSearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func PerformNpmSearch(query string) tea.Cmd {
	return func() tea.Msg {
		if query == "" {
			return SearchResultMsg{Packages: []NpmPackageInfo{}}
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

		var searchResult NpmSearchResponse
		err = json.Unmarshal(body, &searchResult)
		if err != nil {
			return SearchResultMsg{Error: err}
		}

		packages := make([]NpmPackageInfo, 0, len(searchResult.Objects))
		for _, obj := range searchResult.Objects {
			packages = append(packages, obj.Package)
		}

		return SearchResultMsg{Packages: packages}
	}
}

func DebounceSearch(query string, delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(delay)
		return DebouncedSearchMsg{Query: query}
	}
}

func (m NpmSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, HandleSpaceKey(&m)
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
			return m, ToggleAllItems(&m)
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
				return m, tea.Batch(cmd, DebounceSearch(m.Search.Input.Value(), 500*time.Millisecond))
			}
			return m, cmd
		}

	case DebouncedSearchMsg:
		if k.Query != "" && k.Query == m.Search.Input.Value() {
			m.Loading = true
			return m, PerformNpmSearch(k.Query)
		}

	case SearchResultMsg:
		m.Loading = false
		m.Error = k.Error
		m.Packages = k.Packages
		choices := ConvertPackagesToChoices(k.Packages)
		m.Search.Choices = choices
		m.Search.Filtered = choices
		m.Viewport.Cursor = 0
		return m, nil
	}

	return m, nil
}

func HandleSpaceKey(m *NpmSearchModel) tea.Cmd {
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
	return nil
}

func ToggleAllItems(m *NpmSearchModel) tea.Cmd {
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
	return nil
}

func (m NpmSearchModel) View() string {
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

	if m.Loading {
		return BuildStatusLine(&sb, config.Styles.UserGuide.PaddingLeft(1).Render("⌛ searching..."), footerText)
	}

	if m.Error != nil {
		return BuildStatusLine(&sb, config.Styles.Status(config.Colors.Error).PaddingLeft(1).Render(fmt.Sprintf("✗ Error: %v", m.Error)), footerText)
	}

	if len(m.Search.Filtered) == 0 && m.Search.Input.Value() != "" {
		return BuildStatusLine(&sb, config.Styles.UserGuide.PaddingLeft(1).Render("ⓘ  no packages found"), footerText)
	}

	RenderPackageItems(&sb, m)

	sb.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter install"))
	if len(m.Selected) > 0 {
		sb.WriteString(config.Styles.UserGuide.Render(fmt.Sprintf(" (%d selected)", len(m.Selected))))
	}
	sb.WriteString(config.Styles.UserGuide.Render(" • ctrl+a toggle all • esc cancel"))

	return sb.String()
}

func BuildStatusLine(sb *strings.Builder, content string, footer string) string {
	sb.WriteString(config.Styles.Menu.Render("│"))
	sb.WriteString(content)
	sb.WriteString("\n")
	sb.WriteString(config.Styles.UserGuide.Render(footer))
	return sb.String()
}

func RenderPackageItems(sb *strings.Builder, m NpmSearchModel) {
	filtered := len(m.Search.Filtered)
	if filtered == 0 {
		return
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

func InstallNpmPackages(packages []string, bun string) error {
	if len(packages) == 0 {
		return nil
	}

	appDir := "app"
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		return fmt.Errorf("app directory does not exist")
	}

	packageJsonPath := filepath.Join(appDir, "package.json")
	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		return fmt.Errorf("package.json not found in app directory")
	}

	successCount := 0
	for _, pkgName := range packages {
		s := spinner.New(fmt.Sprintf("installing %s to app/node_modules", pkgName))
		go spinner.Start(s)

		// Install the package using bun add
		cmd := exec.Command(bun, "add", pkgName)
		cmd.Dir = appDir
		cmd.Env = os.Environ()

		output, err := cmd.CombinedOutput()
		spinner.Stop(s)

		if err != nil {
			messages.Error(fmt.Sprintf("Failed to install %s: %v\nOutput: %s", pkgName, err, string(output)))
			continue
		}

		messages.Success(fmt.Sprintf("Installed %s to app/node_modules", pkgName))
		successCount++
	}

	if successCount > 0 {
		messages.Success(fmt.Sprintf("Successfully installed %d package(s) to app/node_modules", successCount))
	}

	return nil
}

func Npm(o NpmOptions) error {
	model := InitNpmSearch()

	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run npm search: %w", err)
	}

	m := finalModel.(NpmSearchModel)

	if m.Quitting && !m.Confirmed {
		messages.Info("npm search cancelled")
		return nil
	}

	if len(m.Selected) == 0 {
		messages.Info("No packages selected")
		return nil
	}

	selectedPackages := ExtractPackageNames(m.Selected)

	bunPath := filepath.Join(".gen", "bun", "bun")
	if _, err := os.Stat(bunPath); os.IsNotExist(err) {
		bunPath = "bun"
	}

	return InstallNpmPackages(selectedPackages, bunPath)
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
