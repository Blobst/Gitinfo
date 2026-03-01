package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Repo struct {
	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
}

/* ---------- State Machine ---------- */

type state int

const (
	stateLoading state = iota
	stateTable
	stateError
)

/* ---------- Messages ---------- */

type reposMsg []Repo

type errMsg struct {
	err error
}

/* ---------- Model ---------- */

type model struct {
	table table.Model
	page  int
	user  string
	repos []Repo

	state state
	err   error
}

/* ---------- Bubble Tea Interface ---------- */

func (m model) Init() tea.Cmd {
	return FetchReposCmd(m.user, m.page)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case reposMsg:
		m.state = stateTable
		m.repos = msg
		m.table = BuildTable(msg)
		return m, nil

	case errMsg:
		m.state = stateError
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit
		}

	}

	return m, nil
}

func (m model) View() string {
	switch m.state {

	case stateLoading:
		return "\n Loading repos...\n"

	case stateError:
		return fmt.Sprintf("\n Error: %v\n", m.err)

	default:
		if m.table.Rows() == nil {
			return "\n No repos loaded\n"
		}
		return m.table.View()
	}
}
