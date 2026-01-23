package servers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func Snapshot(server *Server, directory string) (err error) {
	statics := make([]string, 0)
	for _, route := range server.Routes {
		if parts := strings.SplitN(route.Pattern, " ", 2); len(parts) >= 2 {
			if parts[0] != "GET" {
				continue
			}
			if strings.Contains(parts[1], "{") && strings.Contains(parts[1], "}") {
				continue
			}
			statics = append(statics, parts[1])
		}
	}
	var mut sync.Mutex
	errs := make([]error, 0)
	client := http.Client{}
	var group sync.WaitGroup
	for _, static := range statics {
		group.Go(func() {
			mut.Lock()
			defer mut.Unlock()
			url := fmt.Sprintf("http://127.0.0.1:8080%s", static)
			var path string
			if parts := strings.SplitN(strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://"), "/", 2); len(parts) >= 2 {
				path = parts[1]
			}
			if err = os.MkdirAll(filepath.Join(directory, path), os.ModePerm); err != nil {
				errs = append(errs, err)
				return
			}
			fsPath := strings.ReplaceAll(path, "/", string(filepath.Separator))
			var staticResponse *http.Response
			if staticResponse, err = client.Get(url); err != nil {
				errs = append(errs, err)
				return
			}
			defer func(body io.ReadCloser) {
				if cerr := body.Close(); cerr != nil {
					errs = append(errs, cerr)
				}
			}(staticResponse.Body)
			var data []byte
			if data, err = io.ReadAll(staticResponse.Body); err != nil {
				errs = append(errs, err)
				return
			}
			if err = os.WriteFile(filepath.Join(directory, fsPath, "index.html"), data, os.ModePerm); err != nil {
				errs = append(errs, err)
				return
			}
			var request *http.Request
			if request, err = http.NewRequest("GET", url, nil); err != nil {
				errs = append(errs, err)
				return
			}
			request.Header.Add("Accept", "application/json")
			if staticResponse, err = client.Do(request); err != nil {
				return
			}
			defer func(body io.ReadCloser) {
				if cerr := body.Close(); cerr != nil {
					errs = append(errs, err)
				}
			}(staticResponse.Body)
			if data, err = io.ReadAll(staticResponse.Body); err != nil {
				errs = append(errs, err)
				return
			}
			if err = os.WriteFile(filepath.Join(directory, fsPath, "data.json"), data, os.ModePerm); err != nil {
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
