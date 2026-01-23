package snapshots

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func Generate(url string, directory string) (err error) {
	client := http.Client{}
	var staticsResponse *http.Response
	if staticsResponse, err = client.Get(url); err != nil {
		return
	}
	defer func(body io.ReadCloser) {
		if cerr := body.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}(staticsResponse.Body)
	var staticsData []byte
	if staticsData, err = io.ReadAll(staticsResponse.Body); err != nil {
		return
	}
	var statics []string
	if err = json.Unmarshal(staticsData, &statics); err != nil {
		return
	}
	var mut sync.Mutex
	errs := make([]error, 0)
	var group sync.WaitGroup
	for _, static := range statics {
		group.Go(func() {
			mut.Lock()
			defer mut.Unlock()
			routeUrl := fmt.Sprintf("http://127.0.0.1:8080%s", static)
			var path string
			if parts := strings.SplitN(strings.TrimPrefix(strings.TrimPrefix(routeUrl, "https://"), "http://"), "/", 2); len(parts) >= 2 {
				path = parts[1]
			}
			if err = os.MkdirAll(filepath.Join(directory, path), os.ModePerm); err != nil {
				errs = append(errs, err)
				return
			}
			fsPath := strings.ReplaceAll(path, "/", string(filepath.Separator))
			attempts := 0
			attemptsErrs := make([]error, 0)
			var staticResponse *http.Response
			for {
				if staticResponse, err = client.Get(routeUrl); err != nil {
					attemptsErrs = append(attemptsErrs, err)
					if attempts++; attempts >= 2 {
						errs = append(errs, attemptsErrs...)
						break
					}
					time.Sleep(300 * time.Millisecond)
					continue
				}
				break
			}
			defer func(body io.ReadCloser) {
				if cerr := body.Close(); cerr != nil {
					errs = append(errs, cerr)
				}
			}(staticResponse.Body)
			var staticData []byte
			if staticData, err = io.ReadAll(staticResponse.Body); err != nil {
				errs = append(errs, err)
				return
			}
			if err = os.WriteFile(filepath.Join(directory, fsPath, "index.html"), staticData, os.ModePerm); err != nil {
				errs = append(errs, err)
				return
			}
			var request *http.Request
			if request, err = http.NewRequest("GET", routeUrl, nil); err != nil {
				errs = append(errs, err)
				return
			}
			request.Header.Add("Accept", "application/json")
			attemptsErrs = make([]error, 0)
			for {
				if staticResponse, err = client.Do(request); err != nil {
					attemptsErrs = append(attemptsErrs, err)
					if attempts++; attempts >= 2 {
						errs = append(errs, attemptsErrs...)
						break
					}
					time.Sleep(300 * time.Millisecond)
					continue
				}
				break
			}
			defer func(body io.ReadCloser) {
				if cerr := body.Close(); cerr != nil {
					errs = append(errs, err)
				}
			}(staticResponse.Body)
			if staticData, err = io.ReadAll(staticResponse.Body); err != nil {
				errs = append(errs, err)
				return
			}
			if err = os.WriteFile(filepath.Join(directory, fsPath, "data.json"), staticData, os.ModePerm); err != nil {
				errs = append(errs, err)
				return
			}
		})
	}
	group.Wait()
	if len(errs) != 0 {
		err = errors.Join(errs...)
		return
	}
	err = files.CopyDirectory(filepath.Join("app", "dist", "client", "assets"), filepath.Join(directory, "assets"))
	return
}
