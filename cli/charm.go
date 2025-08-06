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

var colors = struct {
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

var styles = struct {
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
	title: lipgloss.NewStyle().Foreground(lipgloss.Color(colors.primary)).Bold(true),
	item:  lipgloss.NewStyle().PaddingLeft(2),
	selected: lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(colors.secondary)),
	status: func(color string) lipgloss.Style {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true).PaddingLeft(2)
	},
	bigText: lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.info)).Bold(true).Align(lipgloss.Center).
		Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color(colors.primary)).
		Padding(1, 4).Margin(1, 2),
	section: lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.secondary)).Bold(true).Underline(true).Padding(1, 0),
	spinner:  lipgloss.NewStyle().Foreground(lipgloss.Color(colors.info)),
	flag:     lipgloss.NewStyle().Foreground(lipgloss.Color(colors.secondary)).Bold(true),
	category: lipgloss.NewStyle().Foreground(lipgloss.Color(colors.primary)).Bold(true).Underline(true),
	example:  lipgloss.NewStyle().Foreground(lipgloss.Color(colors.muted)),
}

func runProgram[T tea.Model](m T) (T, error) {
	result, err := tea.NewProgram(m).Run()
	if err != nil {
		return m, err
	}
	if typed, ok := result.(T); ok {
		return typed, nil
	}
	return m, fmt.Errorf("unexpected model type")
}

func CharmChoose(prompt string, options []string) (string, error) {
	m := simpleChooseModel{choices: options, prompt: prompt}
	result, err := runProgram(m)
	if err != nil {
		return "", err
	}
	if result.selected == "" {
		return "", fmt.Errorf("no selection made")
	}
	return result.selected, nil
}

func CharmMultiSelect(prompt string, options []string) ([]string, error) {
	m := multiSelectModel{choices: options, selected: make(map[int]bool), prompt: prompt}
	result, err := runProgram(m)
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
	m := inputModel{textInput: ti, prompt: prompt}
	result, err := runProgram(m)
	if err != nil {
		return "", err
	}
	return result.textInput.Value(), nil
}

func CharmConfirm(prompt string, defaultValue bool) (bool, error) {
	m := confirmModel{prompt: prompt, defaultValue: defaultValue, confirmed: defaultValue}
	result, err := runProgram(m)
	if err != nil {
		return false, err
	}
	return result.confirmed, nil
}

func CharmSuccess(text string) {
	fmt.Println(styles.status(colors.success).Render("✓ " + text))
}

func CharmError(text string) {
	fmt.Println(styles.status(colors.error).Render("✗ " + text))
}

func CharmWarning(text string) {
	fmt.Println(styles.status(colors.warning).Render("⚠ " + text))
}

func CharmInfo(text string) {
	fmt.Println(styles.status(colors.info).Render("ℹ " + text))
}

func CharmSection(text string) {
	fmt.Println(styles.section.Render("## " + text))
}

type helpSection struct {
	name  string
	flags [][]string
}

var helpData = struct {
	title       string
	description string
	usage       string
	sections    []helpSection
	examples    [][]string
	footer      string
}{
	title:       "Frizzante",
	description: "A modern web framework for Go",
	usage:       "frizzante [OPTIONS]",
	sections: []helpSection{
		{"PROJECT MANAGEMENT", [][]string{
			{"-c, --create-project", "Create a new Frizzante project"},
			{"-a, --add", "Add features to the project (use -a? for details)"},
			{"-i, --install", "Install project dependencies"},
			{"-u, --update", "Update project dependencies"},
			{"    --configure", "Configure project by installing necessary binaries"},
		}},
		{"DEVELOPMENT", [][]string{
			{"-d, --dev", "Start development mode with hot reload"},
			{"-b, --build", "Build project for production"},
			{"-p, --package", "Package app (result in app/dist)"},
			{"    --package-watch", "Watch and package app continuously"},
			{"-t, --test", "Run tests"},
			{"    --check", "Check source code for errors"},
			{"-f, --format", "Format source code"},
		}},
		{"DATABASE", [][]string{
			{"    --sqlc-generate", "Generate SQLC queries"},
		}},
		{"UTILITIES", [][]string{
			{"    --clean", "Clean project build artifacts"},
			{"    --touch", "Create placeholders in app/dist (for go:embed)"},
			{"    --welcome", "Show welcome message"},
		}},
		{"OPTIONS", [][]string{
			{"-h, --help", "Show this help message"},
			{"-v, --version", "Show version information"},
			{"-y, --yes", "Confirm all prompts automatically"},
			{"    --platform", "Set target platform (linux/amd64, darwin/arm64, etc.)"},
			{"    --go", "Set the Go binary path"},
			{"    --air", "Set the Air binary path"},
			{"    --bun", "Set the Bun binary path"},
			{"    --sqlc", "Set the SQLC binary path"},
		}},
	},
	examples: [][]string{
		{"# Create a new project", "frizzante -c my-app"},
		{"# Start development mode", "frizzante -d"},
		{"# Build for production", "frizzante -b"},
		{"# Add features interactively", "frizzante -a:pick"},
	},
	footer: "For more information, visit: https://razshare.github.io/frizzante-docs/",
}

func CharmHelp() {
	CharmBigText(helpData.title)
	fmt.Println(styles.title.Render(helpData.description))
	fmt.Println()

	fmt.Println(styles.category.Render("USAGE:"))
	fmt.Println("  " + helpData.usage)
	fmt.Println()

	for _, section := range helpData.sections {
		fmt.Println(styles.category.Render(section.name + ":"))
		for _, flag := range section.flags {
			fmt.Printf("  %s %s\n", styles.flag.Render(flag[0]), flag[1])
		}
		fmt.Println()
	}

	fmt.Println(styles.category.Render("EXAMPLES:"))
	for _, example := range helpData.examples {
		fmt.Println(styles.example.Render("  " + example[0]))
		fmt.Println("  " + example[1])
		fmt.Println()
	}

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(colors.info)).Italic(true).Render(helpData.footer))
}

func CharmBigText(text string) {
	bigText := generateBigText(text)
	fmt.Println(styles.bigText.Render(bigText))
}

func generateBigText(text string) string {
	if strings.ToLower(text) == "frizzante" {
		return `
███████╗██████╗ ██╗███████╗███████╗ █████╗ ███╗   ██╗████████╗███████╗
██╔════╝██╔══██╗██║╚══███╔╝╚══███╔╝██╔══██╗████╗  ██║╚══██╔══╝██╔════╝
█████╗  ██████╔╝██║  ███╔╝   ███╔╝ ███████║██╔██╗ ██║   ██║   █████╗  
██╔══╝  ██╔══██╗██║ ███╔╝   ███╔╝  ██╔══██║██║╚██╗██║   ██║   ██╔══╝  
██║     ██║  ██║██║███████╗███████╗██║  ██║██║ ╚████║   ██║   ███████╗
╚═╝     ╚═╝  ╚═╝╚═╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝   ╚══════╝
`
	}
	return strings.ToUpper(text)
}

func CharmSpinner(message string) *SpinnerManager {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.spinner
	
	m := spinnerModel{
		spinner: s,
		message: message,
	}
	
	return &SpinnerManager{
		model:   m,
		program: tea.NewProgram(m),
	}
}

func CharmTable(headers []string, rows [][]string) {
	columns := make([]table.Column, len(headers))
	for i, header := range headers {
		width := len(header)
		for _, row := range rows {
			if i < len(row) && len(row[i]) > width {
				width = len(row[i])
			}
		}
		if width > 40 {
			width = 40
		}
		columns[i] = table.Column{Title: header, Width: width + 2}
	}

	tableRows := make([]table.Row, len(rows))
	for i, row := range rows {
		tableRows[i] = row
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(tableRows),
		table.WithFocused(false),
		table.WithHeight(len(rows)+2),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	fmt.Println(t.View())
}

type SpinnerManager struct {
	model   spinnerModel
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

type simpleChooseModel struct {
	choices  []string
	cursor   int
	prompt   string
	selected string
}

func (m simpleChooseModel) Init() tea.Cmd {
	return nil
}

func (m simpleChooseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "enter":
			m.selected = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m simpleChooseModel) View() string {
	s := "\n" + styles.title.Render(m.prompt) + "\n\n"
	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "▶ "
		}
		line := cursor + choice
		if m.cursor == i {
			s += styles.selected.Render(line) + "\n"
		} else {
			s += styles.item.Render(line) + "\n"
		}
	}
	s += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color(colors.muted)).Render("(↑/↓ to navigate, enter to select, q to quit)")
	return s
}

type multiSelectModel struct {
	choices  []string
	cursor   int
	selected map[int]bool
	prompt   string
}

func (m multiSelectModel) Init() tea.Cmd {
	return nil
}

func (m multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m multiSelectModel) View() string {
	s := "\n" + styles.title.Render(m.prompt) + "\n\n"
	s += styles.status(colors.info).Render("Use arrow keys to navigate, space to select, enter to confirm") + "\n\n"
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
			s += styles.selected.Render(line) + "\n"
		} else {
			s += styles.item.Render(line) + "\n"
		}
	}
	return s
}

type inputModel struct {
	textInput textinput.Model
	prompt    string
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m inputModel) View() string {
	return fmt.Sprintf("\n%s\n\n%s\n\n%s",
		styles.title.Render(m.prompt),
		m.textInput.View(),
		"(esc to quit)")
}

type confirmModel struct {
	prompt       string
	confirmed    bool
	defaultValue bool
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m confirmModel) View() string {
	defaultHint := "n"
	if m.defaultValue {
		defaultHint = "Y"
	}
	return "\n" + styles.title.Render(m.prompt) + "\n\n(Y/n) [default: " + defaultHint + "]"
}

type spinnerModel struct {
	spinner spinner.Model
	message string
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m spinnerModel) View() string {
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}