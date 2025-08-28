package npmselect

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli/npm"
	"io"
	"net/http"
	"net/url"
	"time"
)

func Search(query string) tea.Cmd {
	return func() tea.Msg {
		if query == "" {
			return SearchResultMsg{Packages: []npm.PackageInfo{}}
		}

		encodedQuery := url.QueryEscape(query)
		apiUrl := fmt.Sprintf("https://registry.npmjs.org/-/v1/search?text=%s&size=20", encodedQuery)

		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			return SearchResultMsg{Error: err}
		}

		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return SearchResultMsg{Error: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return SearchResultMsg{Error: fmt.Errorf("npm registry returned status %d", resp.StatusCode)}
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return SearchResultMsg{Error: err}
		}

		var searchResult npm.SearchResponse
		err = json.Unmarshal(body, &searchResult)
		if err != nil {
			return SearchResultMsg{Error: err}
		}

		packages := make([]npm.PackageInfo, 0, len(searchResult.Objects))
		for _, obj := range searchResult.Objects {
			packages = append(packages, obj.Package)
		}

		return SearchResultMsg{Packages: packages}
	}
}
