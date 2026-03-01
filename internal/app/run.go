package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Styles ────────────────────────────────────────────────────────────────────

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))

	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

// ── List item ─────────────────────────────────────────────────────────────────

type repoItem struct{ repo Repo }

func (r repoItem) Title() string { return r.repo.Name }
func (r repoItem) Description() string {
	desc := r.repo.Description
	if desc == "" {
		desc = "No description"
	}
	lang := r.repo.Language
	if lang == "" {
		lang = "Unknown"
	}
	return fmt.Sprintf("%s  ⭐ %d  🍴 %d  [%s]", desc, r.repo.Stars, r.repo.Forks, lang)
}
func (r repoItem) FilterValue() string { return r.repo.Name }

// ── Messages ──────────────────────────────────────────────────────────────────

type reposFetchedMsg struct{ repos []Repo }
type errMsg struct{ err error }

// ── State ─────────────────────────────────────────────────────────────────────

type state int

const (
	stateInput state = iota
	stateLoading
	stateList
)

// ── Model ─────────────────────────────────────────────────────────────────────

type Model struct {
	state    state
	input    textinput.Model
	spinner  spinner.Model
	list     list.Model
	username string
	err      error
	width    int
	height   int
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter GitHub username..."
	ti.CharLimit = 64
	ti.Width = 40
	ti.Focus()

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("205")).
		BorderForeground(lipgloss.Color("205"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("240")).
		BorderForeground(lipgloss.Color("205"))

	li := list.New(nil, delegate, 0, 0)
	li.SetShowStatusBar(true)
	li.SetFilteringEnabled(true)
	li.Styles.Title = titleStyle

	return Model{
		state:   stateInput,
		input:   ti,
		spinner: sp,
		list:    li,
	}
}

// Run starts the Bubble Tea program.
func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// ── Commands ──────────────────────────────────────────────────────────────────

func (m Model) fetchRepos() tea.Cmd {
	return func() tea.Msg {
		repos, err := FetchRepos(m.username, 1)
		if err != nil {
			return errMsg{err}
		}
		return reposFetchedMsg{repos}
	}
}

// ── Bubble Tea interface ───────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.state {
		case stateInput:
			if msg.String() == "enter" && m.input.Value() != "" {
				m.username = m.input.Value()
				m.state = stateLoading
				return m, tea.Batch(m.spinner.Tick, m.fetchRepos())
			}

		case stateList:
			if msg.String() == "esc" && !m.list.SettingFilter() {
				// Back to search
				m.state = stateInput
				m.err = nil
				m.input.SetValue("")
				m.input.Focus()
				m.list.SetItems(nil)
				return m, textinput.Blink
			}
		}

	case reposFetchedMsg:
		items := make([]list.Item, len(msg.repos))
		for i, r := range msg.repos {
			items[i] = repoItem{r}
		}
		m.list.SetItems(items)
		m.list.Title = fmt.Sprintf("Repos — %s  (%d)", m.username, len(msg.repos))
		m.state = stateList
		return m, nil

	case errMsg:
		m.err = msg.err
		m.state = stateInput
		return m, nil
	}

	// Delegate remaining updates to the active component
	var cmd tea.Cmd
	switch m.state {
	case stateInput:
		m.input, cmd = m.input.Update(msg)
	case stateLoading:
		m.spinner, cmd = m.spinner.Update(msg)
	case stateList:
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.state {

	case stateInput:
		errText := ""
		if m.err != nil {
			errText = "\n" + errorStyle.Render("✗ "+m.err.Error())
		}
		return docStyle.Render(
			titleStyle.Render("GitHub Repo Viewer") + "\n\n" +
				m.input.View() + errText + "\n\n" +
				subtitleStyle.Render("enter to search • ctrl+c to quit"),
		)

	case stateLoading:
		return docStyle.Render(
			m.spinner.View() + " Fetching repos for " +
				titleStyle.Render(m.username) + "...",
		)

	case stateList:
		return docStyle.Render(
			m.list.View() + "\n" +
				subtitleStyle.Render("/ to filter • esc to search again • ctrl+c to quit"),
		)
	}
	return ""
}
