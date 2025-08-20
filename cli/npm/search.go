package npm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func Search(q string) ([]PackageInfo, error) {
	if q == "" {
		return make([]PackageInfo, 0), nil
	}

	c := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://registry.npmjs.org/-/v1/search?text=%s&size=20", url.QueryEscape(q)), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("npm registry returned status %d", res.StatusCode)
	}

	var d []byte
	d, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var val SearchResponse

	err = json.Unmarshal(d, &val)
	if err != nil {
		return nil, err
	}

	p := make([]PackageInfo, 0, len(val.Objects))
	for _, obj := range val.Objects {
		p = append(p, obj.Package)
	}

	return p, nil
}
