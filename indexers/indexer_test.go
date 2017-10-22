package indexer

import (
	"testing"
)

type MockIndexer struct{}

func (m MockIndexer) GetName() string {
	return "MockIndexer"
}

func (m MockIndexer) GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	return []Torrent{
		Torrent{
			Title:       "Torrent1",
			Description: "Mock Torrent1 torrent",
			Link:        "torrent1.torrent",
			Attributes: map[string]string{
				"seeders": "1",
			},
		},
		Torrent{
			Title:       "Torrent2",
			Description: "Mock Torrent2 torrent",
			Link:        "torrent2.torrent",
			Attributes: map[string]string{
				"seeders": "2",
			},
		},
		Torrent{
			Title:       "Torrent3",
			Description: "Mock Torrent3 torrent",
			Link:        "torrent3.torrent",
			Attributes: map[string]string{
				"seeders": "3",
			},
		},
	}, nil
}

func TestAddIndexer(t *testing.T) {
	AddIndexer(MockIndexer{})
	if len(indexers) != 1 {
		t.Error("Indexer not added to list of indexers")
	}
}
