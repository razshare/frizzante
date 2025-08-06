package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// =============================================================================
// Theme Configuration
// =============================================================================

var Colors = struct {
	primary   string
	secondary string
	success   string
	error     string
	warning   string
	info      string
	muted     string
}{
	primary:   "99",
	secondary: "212",
	success:   "82",
	error:     "196",
	warning:   "214",
	info:      "39",
	muted:     "240",
}

var Styles = struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selected     lipgloss.Style
	status       func(color string) lipgloss.Style
	bigText      lipgloss.Style
	section      lipgloss.Style
	spinner      lipgloss.Style
	flag         lipgloss.Style
	category     lipgloss.Style
	example      lipgloss.Style
}{
	title: lipgloss.NewStyle().Foreground(lipgloss.Color(Colors.primary)).Bold(true),
	item:  lipgloss.NewStyle().PaddingLeft(2),
	selected: lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(Colors.secondary)),
	status: func(color string) lipgloss.Style {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true).PaddingLeft(2)
	},
	bigText: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.info)).Bold(true).Align(lipgloss.Center).
		Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color(Colors.primary)).
		Padding(1, 4).Margin(1, 2),
	section: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.secondary)).Bold(true).Underline(true).Padding(1, 0),
	spinner:  lipgloss.NewStyle().Foreground(lipgloss.Color(Colors.info)),
	flag:     lipgloss.NewStyle().Foreground(lipgloss.Color(Colors.secondary)).Bold(true),
	category: lipgloss.NewStyle().Foreground(lipgloss.Color(Colors.primary)).Bold(true).Underline(true),
	example:  lipgloss.NewStyle().Foreground(lipgloss.Color(Colors.muted)),
}

// =============================================================================
// Core Program Execution
// =============================================================================

func RunProgram[T tea.Model](m T) (T, error) {
	result, err := tea.NewProgram(m).Run()
	if err != nil {
		return m, err
	}
	if typed, ok := result.(T); ok {
		return typed, nil
	}
	return m, fmt.Errorf("unexpected model type")
}

// =============================================================================
// Interactive Input Functions
// =============================================================================

func CharmChoose(prompt string, options []string) (string, error) {
	// Initialize the search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Type to filter..."
	searchInput.Width = 80
	
	m := SimpleChooseModel{
		choices:         options,
		filteredChoices: options,
		prompt:          prompt,
		searchInput:     searchInput,
		maxVisible:      6,
		viewportStart:   0,
		cursor:          0,
		searching:       false,
	}
	
	result, err := RunProgram(m)
	if err != nil {
		return "", err
	}
	if result.selected == "" {
		return "", fmt.Errorf("no selection made")
	}
	return result.selected, nil
}

func CharmMultiSelect(prompt string, options []string) ([]string, error) {
	m := MultiSelectModel{choices: options, selected: make(map[int]bool), prompt: prompt}
	result, err := RunProgram(m)
	if err != nil {
		return nil, err
	}
	var selections []string
	for i, choice := range result.choices {
		if result.selected[i] {
			selections = append(selections, choice)
		}
	}
	return selections, nil
}

func CharmInput(prompt string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()
	ti.Width = 50
	m := InputModel{textInput: ti, prompt: prompt}
	result, err := RunProgram(m)
	if err != nil {
		return "", err
	}
	return result.textInput.Value(), nil
}

func CharmConfirm(prompt string, defaultValue bool) (bool, error) {
	m := ConfirmModel{prompt: prompt, defaultValue: defaultValue, confirmed: defaultValue}
	result, err := RunProgram(m)
	if err != nil {
		return false, err
	}
	return result.confirmed, nil
}

// =============================================================================
// Status & Messaging Functions
// =============================================================================

func CharmSuccess(text string) {
	fmt.Println(Styles.status(Colors.success).Render("✓ " + text))
}

func CharmError(text string) {
	fmt.Println(Styles.status(Colors.error).Render("✗ " + text))
}

func CharmWarning(text string) {
	fmt.Println(Styles.status(Colors.warning).Render("⚠ " + text))
}

func CharmInfo(text string) {
	fmt.Println(Styles.status(Colors.info).Render("ℹ " + text))
}

func CharmSection(text string) {
	fmt.Println(Styles.section.Render("## " + text))
}

func CharmDockerHelp() {
	fmt.Println()
	fmt.Println(Styles.title.Render("🐙 You're running Frizzante in Docker!"))
	fmt.Println()
	
	fmt.Println(Styles.section.Render("⚡ Simple workflow:"))
	
	fmt.Println(Styles.status(Colors.info).Render("• Attach to the container: ") + 
		Styles.example.Render("docker exec -it frizzante-start sh"))
	
	fmt.Println(Styles.status(Colors.info).Render("• Run environment in container:"))
	fmt.Println(Styles.item.Render("    • Dev environment: ") + Styles.flag.Render("make dev"))
	fmt.Println(Styles.item.Render("    • Prod environment: ") + Styles.flag.Render("make build"))
	fmt.Println(Styles.item.Render("    • To run the app: ") + Styles.flag.Render("./.gen/bin/app"))
	
	fmt.Println(Styles.status(Colors.info).Render("• Run prod via docker:"))
	fmt.Println(Styles.item.Render("    • Build image: ") + 
		Styles.example.Render("docker build --target frizzante_prod -t my-app:prod ."))
	fmt.Println(Styles.item.Render("    • Run image: ") + 
		Styles.example.Render("docker run -p 8080:8080 my-app:prod"))
	fmt.Println(Styles.item.Render("    • Via docker compose: ") + 
		Styles.example.Render("docker compose -f compose.yaml -f compose.prod.yaml up -d --build"))
	
	fmt.Println(Styles.status(Colors.success).Render("🎉 Enjoy!!"))
	fmt.Println()
}

// =============================================================================
// Display & Formatting Functions
// =============================================================================

func wrapText(text string, width int) []string {
	if width <= 0 || len(text) <= width {
		return []string{text}
	}
	
	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}
	
	currentLine := ""
	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	
	return lines
}

func CharmTable(headers []string, rows [][]string) {
	columns := make([]table.Column, len(headers))
	maxColWidth := 60
	
	for i, header := range headers {
		width := len(header)
		for _, row := range rows {
			if i < len(row) && len(row[i]) > width {
				width = len(row[i])
			}
		}
		if width > maxColWidth {
			width = maxColWidth
		}
		columns[i] = table.Column{Title: header, Width: width + 2}
	}

	wrappedRows := []table.Row{}
	for rowIdx, row := range rows {
		maxLines := 1
		wrappedCells := make([][]string, len(row))
		
		for i, cell := range row {
			if i < len(columns) {
				cellWidth := columns[i].Width - 2
				wrappedCells[i] = wrapText(cell, cellWidth)
				if len(wrappedCells[i]) > maxLines {
					maxLines = len(wrappedCells[i])
				}
			}
		}
		
		for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
			newRow := make([]string, len(row))
			for cellIdx, wrappedCell := range wrappedCells {
				if lineIdx < len(wrappedCell) {
					newRow[cellIdx] = wrappedCell[lineIdx]
				} else {
					newRow[cellIdx] = ""
				}
			}
			wrappedRows = append(wrappedRows, newRow)
		}
		
		// Add empty row after each logical row (except the last one)
		if rowIdx < len(rows)-1 {
			emptyRow := make([]string, len(row))
			for i := range emptyRow {
				emptyRow[i] = ""
			}
			wrappedRows = append(wrappedRows, emptyRow)
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(wrappedRows),
		table.WithFocused(false),
		table.WithHeight(len(wrappedRows)+2),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	// Remove selection highlight
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)

	fmt.Println(t.View())
}

func CharmSpinner(message string) *SpinnerManager {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = Styles.spinner
	
	m := SpinnerModel{
		spinner: s,
		message: message,
	}
	
	return &SpinnerManager{
		Model:   m,
		program: tea.NewProgram(m),
	}
}

// =============================================================================
// Spinner Management
// =============================================================================

type SpinnerManager struct {
	Model   SpinnerModel
	program *tea.Program
	done    chan bool
}

func (sm *SpinnerManager) Start() {
	sm.done = make(chan bool)
	go func() {
		sm.program.Run()
		close(sm.done)
	}()
	time.Sleep(100 * time.Millisecond)
}

func (sm *SpinnerManager) Stop() {
	sm.program.Quit()
	<-sm.done
	fmt.Print("\r\033[K")
}

// Spinner Model - Loading spinner display
type SpinnerModel struct {
	spinner spinner.Model
	message string
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m SpinnerModel) View() string {
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

// =============================================================================
// UI Models (Bubble Tea Components)
// =============================================================================

// Simple Choice Model - Single selection from a list with search and scrolling
type SimpleChooseModel struct {
	choices        []string          // All available choices
	filteredChoices []string         // Choices after filtering
	cursor         int               // Current cursor position in filtered list
	prompt         string            // The prompt to display
	selected       string            // The selected choice
	searchInput    textinput.Model   // Search input field
	searching      bool              // Whether we're in search mode
	viewportStart  int               // Start index for viewport (for scrolling)
	maxVisible     int               // Maximum visible items (5 by default)
}

func (m SimpleChooseModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SimpleChooseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle search mode
		if m.searching {
			switch msg.String() {
			case "esc":
				m.searching = false
				m.searchInput.SetValue("")
				m.filteredChoices = m.choices
				m.cursor = 0
				m.viewportStart = 0
				return m, nil
			case "enter":
				if len(m.filteredChoices) > 0 {
					m.selected = m.filteredChoices[m.cursor]
					return m, tea.Quit
				}
			case "up", "ctrl+p":
				m.navigateUp()
				return m, nil
			case "down", "ctrl+n", "tab":
				m.navigateDown()
				return m, nil
			default:
				// Update search input
				prevValue := m.searchInput.Value()
				m.searchInput, cmd = m.searchInput.Update(msg)
				if m.searchInput.Value() != prevValue {
					m.filterChoices()
				}
				return m, cmd
			}
		}
		
		// Normal navigation mode
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "/", "ctrl+f":
			m.searching = true
			m.searchInput.Focus()
			return m, textinput.Blink
		case "up", "k":
			m.navigateUp()
		case "down", "j":
			m.navigateDown()
		case "g":
			// Go to top
			m.cursor = 0
			m.viewportStart = 0
		case "G":
			// Go to bottom
			if len(m.filteredChoices) > 0 {
				m.cursor = len(m.filteredChoices) - 1
				m.updateViewport()
			}
		case "enter":
			if len(m.filteredChoices) > 0 {
				m.selected = m.filteredChoices[m.cursor]
				return m, tea.Quit
			}
		}
	}
	
	return m, nil
}

func (m *SimpleChooseModel) navigateUp() {
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.viewportStart {
			m.viewportStart = m.cursor
		}
	}
}

func (m *SimpleChooseModel) navigateDown() {
	if m.cursor < len(m.filteredChoices)-1 {
		m.cursor++
		if m.cursor >= m.viewportStart+m.maxVisible {
			m.viewportStart = m.cursor - m.maxVisible + 1
		}
	}
}

func (m *SimpleChooseModel) updateViewport() {
	// Ensure viewport shows the cursor
	if m.cursor < m.viewportStart {
		m.viewportStart = m.cursor
	} else if m.cursor >= m.viewportStart+m.maxVisible {
		m.viewportStart = m.cursor - m.maxVisible + 1
	}
	
	// Ensure viewport doesn't go out of bounds
	if m.viewportStart < 0 {
		m.viewportStart = 0
	}
	maxStart := len(m.filteredChoices) - m.maxVisible
	if maxStart < 0 {
		maxStart = 0
	}
	if m.viewportStart > maxStart {
		m.viewportStart = maxStart
	}
}

func (m *SimpleChooseModel) filterChoices() {
	searchTerm := strings.ToLower(m.searchInput.Value())
	if searchTerm == "" {
		m.filteredChoices = m.choices
	} else {
		m.filteredChoices = []string{}
		for _, choice := range m.choices {
			if strings.Contains(strings.ToLower(choice), searchTerm) {
				m.filteredChoices = append(m.filteredChoices, choice)
			}
		}
	}
	
	// Reset cursor and viewport
	m.cursor = 0
	m.viewportStart = 0
}

func (m SimpleChooseModel) View() string {
	var s strings.Builder
	
	s.WriteString(Styles.title.Render(m.prompt) + "\n")
	
	// Search bar
	if m.searching {
		s.WriteString(Styles.status(Colors.info).Render("🔍 Search: "))
		s.WriteString(m.searchInput.View())
		s.WriteString("\n")
	}
	
	// Show filtered results count if searching
	if m.searching && m.searchInput.Value() != "" {
		s.WriteString(Styles.status(Colors.muted).Render(
			fmt.Sprintf("Found %d results", len(m.filteredChoices))))
		s.WriteString("\n")
	}
	
	// Display choices (only show maxVisible items)
	viewportEnd := m.viewportStart + m.maxVisible
	if viewportEnd > len(m.filteredChoices) {
		viewportEnd = len(m.filteredChoices)
	}
	
	// Show scroll indicator at top
	if m.viewportStart > 0 {
		s.WriteString(Styles.status(Colors.muted).Render("    ↑ more above") + "\n")
	}
	
	for i := m.viewportStart; i < viewportEnd; i++ {
		choice := m.filteredChoices[i]
		cursor := "  "
		if i == m.cursor {
			cursor = "▶ "
		}
		
		line := cursor + choice
		if i == m.cursor {
			s.WriteString(Styles.selected.Render(line) + "\n")
		} else {
			s.WriteString(Styles.item.Render(line) + "\n")
		}
	}
	
	// Show scroll indicator at bottom
	if viewportEnd < len(m.filteredChoices) {
		s.WriteString(Styles.status(Colors.muted).Render("    ↓ more below") + "\n")
	}
	
	// Help text
	if m.searching {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(Colors.muted)).Render(
			"(↑/↓ navigate, enter to select, esc to clear search)"))
	} else {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(Colors.muted)).Render(
			"(↑/↓ navigate, / to search, enter to select, ctrl+c to quit)"))
	}
	
	return s.String()
}

// Multi Select Model - Multiple selections from a list
type MultiSelectModel struct {
	choices  []string
	cursor   int
	selected map[int]bool
	prompt   string
}

func (m MultiSelectModel) Init() tea.Cmd {
	return nil
}

func (m MultiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case " ":
			if m.selected[m.cursor] {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = true
			}
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MultiSelectModel) View() string {
	s := Styles.title.Render(m.prompt) + "\n"
	s += Styles.status(Colors.info).Render("Use arrow keys to navigate, space to select, enter to confirm") + "\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if m.selected[i] {
			checked = "✓"
		}
		line := cursor + " [" + checked + "] " + choice
		if m.cursor == i {
			s += Styles.selected.Render(line) + "\n"
		} else {
			s += Styles.item.Render(line) + "\n"
		}
	}
	return s
}

// Input Model - Text input field
type InputModel struct {
	textInput textinput.Model
	prompt    string
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m InputModel) View() string {
	return fmt.Sprintf("\n%s\n\n%s\n\n%s",
		Styles.title.Render(m.prompt),
		m.textInput.View(),
		"(esc to quit)")
}

// Confirm Model - Yes/No confirmation dialog
type ConfirmModel struct {
	prompt       string
	confirmed    bool
	defaultValue bool
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.confirmed = true
			return m, tea.Quit
		case "n", "N":
			m.confirmed = false
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.confirmed = false
			return m, tea.Quit
		case "enter":
			m.confirmed = m.defaultValue
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ConfirmModel) View() string {
	defaultHint := "n"
	if m.defaultValue {
		defaultHint = "Y"
	}
	return "\n" + Styles.title.Render(m.prompt) + "\n(Y/n) [default: " + defaultHint + "]"
}