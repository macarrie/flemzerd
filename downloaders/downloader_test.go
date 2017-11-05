package downloader

type MockDownloader struct{}

var testTorrentsCount int

func (d MockDownloader) AddTorrent(url string) error {
	testTorrentsCount += 1
	return nil
}

func (d MockDownloader) Init() error {
	return nil
}
