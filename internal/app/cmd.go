package app

import tea "github.com/charmbracelet/bubbletea"

func FetchReposCmd(user string, page int) tea.Cmd {
	return func() tea.Msg {
		repos, err := FetchRepos(user, page)
		if err != nil {
			return errMsg{err}
		}

		return reposMsg(repos)
	}
}
