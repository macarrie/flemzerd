package indexer

type Indexer interface {
	GetName() string
	GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error)
}

type Torrent struct {
	Title       string
	Description string
	Link        string
	Attributes  map[string]string
}
