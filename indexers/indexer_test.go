package indexer

import (
	"errors"
)

type MockIndexer struct{}

func (m MockIndexer) GetName() string {
	return "MockIndexer"
}

func (m MockIndexer) GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	if season == 0 {
		return []Torrent{}, errors.New("Could not get torrents for episode")
	} else if episode == 0 {
		return []Torrent{}, nil
	}

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
