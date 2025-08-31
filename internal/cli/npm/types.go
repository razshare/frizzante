package npm

type PackageInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type SearchResponse struct {
	Objects []struct {
		Package PackageInfo `json:"package"`
	} `json:"objects"`
}

type SearchChannels struct {
	Packages chan []PackageInfo
	Query    chan string
	Error    chan error
	Stop     chan any
}
