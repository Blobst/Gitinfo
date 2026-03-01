package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func FetchRepos(user string, page int) ([]Repo, error) {
	apiURL := fmt.Sprintf(
		"https://api.github.com/users/%s/repos?page=%d&per_page=20",
		url.PathEscape(user),
		page,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var Repos []Repo
	err = json.NewDecoder(resp.Body).Decode(&Repos)

	return Repos, err
}
