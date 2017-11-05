package indexer

import (
	"strconv"
	"testing"
)

func TestAddIndexer(t *testing.T) {
	AddIndexer(MockIndexer{})
	if len(indexersCollection) != 1 {
		t.Error("Indexer not added to list of indexers")
	}
}

func TestGetTorrentForEpisode(t *testing.T) {
	ind1 := MockIndexer{}
	ind2 := MockIndexer{}

	indexersCollection = []Indexer{ind1, ind2}

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
