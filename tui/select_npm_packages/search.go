package select_npm_packages

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli/npm"
)

func Search(query string) tea.Cmd {
	return func() (message tea.Msg) {
		if query == "" {
			message = SearchResultMsg{Packages: []npm.PackageInfo{}}
			return
		}
		encodedQuery := url.QueryEscape(query)
		apiUrl := fmt.Sprintf("https://registry.npmjs.org/-/v1/search?text=%s&size=20", encodedQuery)
		client := &http.Client{
			Timeout: 5 * time.Second,
		}
		request, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			message = SearchResultMsg{Error: err}
			return
		}
		request.Header.Set("Accept", "application/json")
		response, err := client.Do(request)
		if err != nil {
			message = SearchResultMsg{Error: err}
			return
		}
		if response.Body != nil {
			defer func() {
				if cerr := response.Body.Close(); cerr != nil {
					if message == nil {
						message = SearchResultMsg{Error: cerr}
						return
					}
				}
			}()
		}
		if response.StatusCode != http.StatusOK {
			message = SearchResultMsg{Error: fmt.Errorf("npm registry returned status %d", response.StatusCode)}
			return
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			message = SearchResultMsg{Error: err}
			return
		}
		var searchResult npm.SearchResponse
		err = json.Unmarshal(body, &searchResult)
		if err != nil {
			message = SearchResultMsg{Error: err}
			return
		}
		packages := make([]npm.PackageInfo, 0, len(searchResult.Objects))
		for _, obj := range searchResult.Objects {
			packages = append(packages, obj.Package)
		}
		message = SearchResultMsg{Packages: packages}
		return
	}
}
