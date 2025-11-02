package npm

type SearchResponse struct {
	Objects []struct {
		Package PackageInfo `json:"package"`
	} `json:"objects"`
}
