package app

import "github.com/charmbracelet/bubbles/table"

func BuildTable(repos []Repo) table.Model {
	Columns := []table.Column{
		{Title: "HTTP URL", Width: 60},
		{Title: "SSH URL", Width: 50},
	}

	var Rows []table.Row

	for _, r := range repos {
		Rows = append(Rows, table.Row{
			r.CloneURL,
			r.SSHURL,
		})
	}

	t := table.New(
		table.WithColumns(Columns),
		table.WithRows(Rows),
		table.WithFocused(true),
	)

	return t
}
