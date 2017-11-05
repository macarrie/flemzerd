package downloader

type Downloader interface {
	AddTorrent(url string) error
	Init() error
}
