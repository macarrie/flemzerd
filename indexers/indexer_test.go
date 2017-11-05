package indexer

import (
	"errors"
	"strconv"
	"testing"
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

func TestAddIndexer(t *testing.T) {
	AddIndexer(MockIndexer{})
	if len(indexers) != 1 {
		t.Error("Indexer not added to list of indexers")
	}
}

func TestGetTorrentForEpisode(t *testing.T) {
	ind1 := MockIndexer{}
	ind2 := MockIndexer{}

	indexers = []Indexer{ind1, ind2}

	torrentList, _ := GetTorrentForEpisode("Test show", 1, 1)
	if len(torrentList) != 6 {
		t.Errorf("Expected 6 torrents, got %d instead\n", len(torrentList))
	}

	length, _ := strconv.Atoi(torrentList[0].Attributes["seeders"])
	if length != 3 {
		t.Error("Torrent list is not sorted by seeders")
	}

	torrentList, err := GetTorrentForEpisode("Test show", 0, 1)
	if len(torrentList) != 0 || err != nil {
		t.Errorf("Expected errors during indexer search and 0 results, got %d results instead\n", len(torrentList))
	}

	torrentList, err = GetTorrentForEpisode("Test show", 1, 0)
	if err != nil {
		t.Error("Expected to have zero results and no error, go an error instead: ", err)
	}
	if len(torrentList) != 0 {
		t.Errorf("Expected to have no results, got %d results instead\n", len(torrentList))
	}
}
