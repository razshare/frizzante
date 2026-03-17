package generations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/messages"
)

func Snapshot(options SnapshotOptions) (err error) {
	var staticsUrl string
	if staticsUrl = options.StaticsUrl; staticsUrl == "" {
		if options.Strict {
			err = errors.New("statics url is empty")
			return
		}
		if staticsUrl, err = inputs.Send("statics url"); err != nil {
			return
		}
	}
	directoryName := filepath.Join(".gen", "snapshot")
	client := http.Client{}
	var staticsResponse *http.Response
	if staticsResponse, err = client.Get(staticsUrl); err != nil {
		return
	}
	if staticsResponse.Body != nil {
		defer func() {
			if cerr := staticsResponse.Body.Close(); cerr != nil {
				if err == nil {
					err = cerr
				}
			}
		}()
	}
	var staticsData []byte
	if staticsData, err = io.ReadAll(staticsResponse.Body); err != nil {
		return
	}
	var staticPaths []string
	if err = json.Unmarshal(staticsData, &staticPaths); err != nil {
		return
	}
	generate := func(staticPath string) (err error) {
		url := fmt.Sprintf("http://127.0.0.1:8080%s", staticPath)
		var path string
		if parts := strings.SplitN(strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://"), "/", 2); len(parts) >= 2 {
			path = strings.ReplaceAll(parts[1], "/", string(filepath.Separator))
		}
		if err = os.MkdirAll(filepath.Join(directoryName, path), os.ModePerm); err != nil {
			return
		}
		var request *http.Request
		if request, err = http.NewRequest("GET", url, nil); err != nil {
			return
		}
		request.Header.Add("Accept", "text/html")
		request.Header.Add("X-FrizzanteViewType", "snapshot")
		var staticResponseHtml *http.Response
		if staticResponseHtml, err = client.Do(request); err != nil {
			return
		}
		if staticResponseHtml.Body != nil {
			defer func() {
				if cerr := staticResponseHtml.Body.Close(); cerr != nil {
					if err == nil {
						err = cerr
					}
				}
			}()
		}
		var staticData []byte
		if staticData, err = io.ReadAll(staticResponseHtml.Body); err != nil {
			return
		}
		fileNameHtml := filepath.Join(directoryName, path, "index.html")
		if err = os.WriteFile(fileNameHtml, staticData, os.ModePerm); err != nil {
			return
		}
		messages.Successf("%s generated from %s", fileNameHtml, url)
		if request, err = http.NewRequest("GET", url, nil); err != nil {
			return
		}
		request.Header.Add("Accept", "application/json")
		request.Header.Add("X-FrizzanteViewType", "snapshot")
		var staticResponseJson *http.Response
		if staticResponseJson, err = client.Do(request); err != nil {
			return
		}
		if staticResponseJson.Body != nil {
			defer func() {
				if cerr := staticResponseJson.Body.Close(); cerr != nil {
					if err == nil {
						err = cerr
					}
				}
			}()
		}
		var staticDataJson []byte
		if staticDataJson, err = io.ReadAll(staticResponseJson.Body); err != nil {
			return
		}
		fileNameJson := filepath.Join(directoryName, path, "data.json")
		if err = os.WriteFile(fileNameJson, staticDataJson, os.ModePerm); err != nil {
			return
		}
		messages.Successf("%s generated from %s", fileNameJson, url)
		return
	}
	var staticsPath string
	if parts := strings.SplitN(strings.TrimPrefix(strings.TrimPrefix(staticsUrl, "https://"), "http://"), "/", 2); len(parts) >= 2 {
		staticsPath = "/" + parts[1]
	}
	for _, staticPath := range staticPaths {
		if staticPath == staticsPath {
			messages.Infof("skipping %s because it matches the statics path", staticPath)
			continue
		}
		if err = generate(staticPath); err != nil {
			return
		}
	}
	err = files.CopyDirectory(filepath.Join("app", "dist", "client", "assets"), filepath.Join(directoryName, "assets"))
	return
}
