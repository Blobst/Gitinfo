package app

import tea "github.com/charmbracelet/bubbletea"

func Run(user string) {
	m := model{
		page:  1,
		user:  user,
		state: stateLoading,
	}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
