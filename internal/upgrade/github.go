package upgrade

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Release represents a Github release
type Release struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Draft   bool    `json:"draft"`
	Assets  []Asset `json:"assets"`
}

// Asset represents a release asset
type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int    `json:"size"`
}

// GithubClient handles Github API interactions
type GithubClient struct {
	token string
}

// NewGithubClient creates a new Github API client
func NewGithubClient(token string) *GithubClient {
	return &GithubClient{token: token}
}

// GetLatestRelease fetches the latest release
func (gc *GithubClient) GetLatestRelease(owner, repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	return gc.getRelease(url)
}

// GetRelease fetches a specific release by tag
func (gc *GithubClient) GetRelease(owner, repo, tag string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag)
	return gc.getRelease(url)
}

func (gc *GithubClient) getRelease(url string) (*Release, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if gc.token != "" {
		req.Header.Set("Authorization", "token "+gc.token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d: %s", resp.StatusCode, string(body))
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}
