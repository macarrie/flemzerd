package indexer

import (
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

	if torrentList[0].Seeders != 3 {
		t.Error("Torrent list is not sorted by seeders")
	}

	torrentList, err := GetTorrentForEpisode("Test show", 1, 0)
	if err == nil {
		t.Error("Expected to have zero results and an error, go no error instead: ")
	}
	if len(torrentList) != 0 {
		t.Errorf("Expected to have no results, got %d results instead\n", len(torrentList))
	}
}
