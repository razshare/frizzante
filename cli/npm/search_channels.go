package npm

type SearchChannels struct {
	Packages chan []PackageInfo
	Query    chan string
	Error    chan error
	Stop     chan struct{}
}
