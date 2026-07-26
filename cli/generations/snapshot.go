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

	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/tui/inputs"
	"github.com/razshare/frizzante/v2/tui/messages"
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
	var staticsRequest *http.Request
	if staticsRequest, err = http.NewRequest("GET", staticsUrl, nil); err != nil {
		return
	}
	staticsRequest.Header.Add("Accept", "application/json")
	var staticsResponse *http.Response
	if staticsResponse, err = client.Do(staticsRequest); err != nil {
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
	if err = os.RemoveAll(directoryName); err != nil {
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
		var endpointRequest *http.Request
		if endpointRequest, err = http.NewRequest("GET", url, nil); err != nil {
			return
		}
		endpointRequest.Header.Add("Accept", "text/html")
		endpointRequest.Header.Add("X-FrizzanteViewType", "snapshot")
		var staticResponseHtml *http.Response
		if staticResponseHtml, err = client.Do(endpointRequest); err != nil {
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
		if endpointRequest, err = http.NewRequest("GET", url, nil); err != nil {
			return
		}
		endpointRequest.Header.Add("Accept", "application/json")
		endpointRequest.Header.Add("X-FrizzanteViewType", "snapshot")
		var staticResponseJson *http.Response
		if staticResponseJson, err = client.Do(endpointRequest); err != nil {
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
			if staticsPath == "/favicon.ico" {
				messages.Warning("could not snapshot /favicon.ico")
				continue
			}
			return
		}
	}
	err = files.CopyDirectory(filepath.Join("app", "dist", "client", "assets"), filepath.Join(directoryName, "assets"))
	return
}
