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
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

type NpmPackageInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type NpmSearchResponse struct {
	Objects []struct {
		Package NpmPackageInfo `json:"package"`
	} `json:"objects"`
}

type NpmSearchModel struct {
	Input         textinput.Model
	Packages      []NpmPackageInfo
	Selected      map[int]bool
	Cursor        int
	Loading       bool
	Error         error
	LastQuery     string
	DebounceTimer *time.Timer
	Quitting      bool
	Confirmed     bool
}

type SearchResultMsg struct {
	Packages []NpmPackageInfo
	Error    error
}

type DebouncedSearchMsg struct {
	Query string
}

func InitNpmSearch() NpmSearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search npm packages..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	return NpmSearchModel{
		Input:    ti,
		Packages: []NpmPackageInfo{},
		Selected: make(map[int]bool),
		Cursor:   0,
		Loading:  false,
	}
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
		apiUrl := fmt.Sprintf("https://registry.npmjs.org/-/v1/search?text=%s&size=15", encodedQuery)

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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit

		case "enter":
			if len(m.Selected) > 0 {
				m.Confirmed = true
				return m, tea.Quit
			}
			return m, nil

		case "up", "ctrl+p":
			if len(m.Packages) > 0 && m.Cursor > 0 {
				m.Cursor--
			}
			return m, nil

		case "down", "ctrl+n":
			if len(m.Packages) > 0 && m.Cursor < len(m.Packages)-1 {
				m.Cursor++
			}
			return m, nil

		case " ":
			if len(m.Packages) > 0 && m.Cursor < len(m.Packages) {
				if m.Selected[m.Cursor] {
					delete(m.Selected, m.Cursor)
				} else {
					m.Selected[m.Cursor] = true
				}
				return m, nil
			}
		
		case "ctrl+a":
			if len(m.Packages) > 0 {
				if len(m.Selected) == len(m.Packages) {
					m.Selected = make(map[int]bool)
				} else {
					for i := range m.Packages {
						m.Selected[i] = true
					}
				}
				return m, nil
			}
		}
		
		// For all other keys, update the input
		var cmd tea.Cmd
		prevValue := m.Input.Value()
		m.Input, cmd = m.Input.Update(msg)
		
		if m.Input.Value() != prevValue {
			m.LastQuery = m.Input.Value()
			if m.DebounceTimer != nil {
				m.DebounceTimer.Stop()
			}
			return m, tea.Batch(cmd, DebounceSearch(m.Input.Value(), 500*time.Millisecond))
		}
		return m, cmd

	case DebouncedSearchMsg:
		if msg.Query != "" && msg.Query == m.Input.Value() {
			m.Loading = true
			return m, PerformNpmSearch(msg.Query)
		}

	case SearchResultMsg:
		m.Loading = false
		m.Error = msg.Error
		m.Packages = msg.Packages
		m.Cursor = 0
		m.Selected = make(map[int]bool)
		return m, nil
	}

	return m, nil
}

func (m NpmSearchModel) View() string {
	var sb strings.Builder
	sb.Grow(1024)

	sb.WriteString(config.Styles.Menu.Render("Search NPM packages"))

	if m.Input.Value() != "" {
		sb.WriteString(config.Styles.UserInput.Render(" ⁋/" + m.Input.Value()))
	} else {
		sb.WriteString(config.Styles.UserGuide.Render(" ⁋/type to search"))
	}

	sb.WriteString("\n")

	if m.Loading {
		sb.WriteString(config.Styles.Menu.Render("│"))
		sb.WriteString(config.Styles.UserGuide.PaddingLeft(1).Render("⌛ searching..."))
		sb.WriteString("\n")
		sb.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter install • ctrl+a toggle all • esc cancel"))
		return sb.String()
	}

	if m.Error != nil {
		sb.WriteString(config.Styles.Menu.Render("│"))
		sb.WriteString(config.Styles.Status(config.Colors.Error).PaddingLeft(1).Render(fmt.Sprintf("✗ Error: %v", m.Error)))
		sb.WriteString("\n")
		sb.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter install • ctrl+a toggle all • esc cancel"))
		return sb.String()
	}

	if len(m.Packages) == 0 && m.Input.Value() != "" {
		sb.WriteString(config.Styles.Menu.Render("│"))
		sb.WriteString(config.Styles.UserGuide.PaddingLeft(1).Render("ⓘ  no packages found"))
		sb.WriteString("\n")
		sb.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter install • ctrl+a toggle all • esc cancel"))
		return sb.String()
	}

	for i, pkg := range m.Packages {
		sb.WriteString(config.Styles.Menu.Render("│"))
		
		pkgName := pkg.Name
		if pkg.Version != "" {
			pkgName = fmt.Sprintf("%s@%s", pkg.Name, pkg.Version)
		}

		if i == m.Cursor {
			if m.Selected[i] {
				sb.WriteString(config.Styles.Selected.Render("● " + pkgName))
			} else {
				sb.WriteString(config.Styles.Selected.Render("◉ " + pkgName))
			}
			
			if pkg.Description != "" {
				desc := pkg.Description
				if len(desc) > 50 {
					desc = desc[:47] + "..."
				}
				sb.WriteString(config.Styles.UserGuide.Render("  ⇢  " + desc))
			}
		} else if m.Selected[i] {
			sb.WriteString(config.Styles.Item.Render("● " + pkgName))
		} else {
			sb.WriteString(config.Styles.Item.Render("○ " + pkgName))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter install"))
	if len(m.Selected) > 0 {
		sb.WriteString(config.Styles.UserGuide.Render(fmt.Sprintf(" (%d selected)", len(m.Selected))))
	}
	sb.WriteString(config.Styles.UserGuide.Render(" • ctrl+a toggle all • esc cancel"))

	return sb.String()
}

func InstallNpmPackages(packages []string, bun string) error {
	if len(packages) == 0 {
		return nil
	}

	// Create .gen directory if it doesn't exist
	err := os.MkdirAll(".gen", 0755)
	if err != nil {
		return fmt.Errorf("failed to create .gen directory: %w", err)
	}

	// NOT SURE IF WE WANT TO DO ANY OF THE BELOW, I GUESS WE WOULD JUST INSTALL IN APP/node_modules
	successCount := 0
	for _, pkgName := range packages {
		s := spinner.New(fmt.Sprintf("installing %s to .gen/%s", pkgName, pkgName))
		go spinner.Start(s)

		// Create directory for this package
		pkgDir := filepath.Join(".gen", pkgName)
		err := os.MkdirAll(pkgDir, 0755)
		if err != nil {
			spinner.Stop(s)
			messages.Error(fmt.Sprintf("Failed to create directory for %s: %v", pkgName, err))
			continue
		}

		// Create a minimal package.json for this package
		packageJson := fmt.Sprintf(`{
  "name": "%s-wrapper",
  "version": "1.0.0",
  "private": true,
  "dependencies": {
    "%s": "latest"
  }
}`, pkgName, pkgName)

		packageJsonPath := filepath.Join(pkgDir, "package.json")
		err = os.WriteFile(packageJsonPath, []byte(packageJson), 0644)
		if err != nil {
			spinner.Stop(s)
			messages.Error(fmt.Sprintf("Failed to create package.json for %s: %v", pkgName, err))
			continue
		}

		// Install the package using bun
		cmd := exec.Command(bun, "install")
		cmd.Dir = pkgDir
		cmd.Env = append(os.Environ())
		
		output, err := cmd.CombinedOutput()
		spinner.Stop(s)

		if err != nil {
			messages.Error(fmt.Sprintf("Failed to install %s: %v\nOutput: %s", pkgName, err, string(output)))
			continue
		}

		messages.Success(fmt.Sprintf("Installed %s to .gen/%s", pkgName, pkgName))
		successCount++
	}

	if successCount > 0 {
		messages.Success(fmt.Sprintf("Successfully installed %d package(s) to .gen", successCount))
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

	selectedPackages := []string{}
	for idx := range m.Selected {
		if idx < len(m.Packages) {
			pkg := m.Packages[idx]
			selectedPackages = append(selectedPackages, pkg.Name)
		}
	}

	bunPath := filepath.Join(".gen", "bun", "bun")
	if _, err := os.Stat(bunPath); os.IsNotExist(err) {
		bunPath = "bun"
	}

	return InstallNpmPackages(selectedPackages, bunPath)
}