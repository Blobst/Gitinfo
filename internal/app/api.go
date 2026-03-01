package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	SSHURL      string `json:"ssh_url"`
	CloneURL    string `json:"clone_url"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Private     bool   `json:"private"`
}

// SSHCloneURL returns the SSH clone URL, e.g. git@github.com:user/repo.git
func (r Repo) SSHCloneURL() string {
	return r.SSHURL
}

// HTTPCloneURL returns the HTTPS clone URL, e.g. https://github.com/user/repo.git
func (r Repo) HTTPCloneURL() string {
	return r.CloneURL
}

func FetchRepos(user string, page int) ([]Repo, error) {
	apiURL := fmt.Sprintf(
		"https://api.github.com/users/%s/repos?page=%d&per_page=20",
		url.PathEscape(user),
		page,
	)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user %q not found", user)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API error: %s", resp.Status)
	}

	var repos []Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}
	return repos, nil
}
