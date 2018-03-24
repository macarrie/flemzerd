package indexer

import (
	"testing"
)

func TestAddIndexer(t *testing.T) {
	AddIndexer(MockTVIndexer{})
	if len(indexersCollection) != 1 {
		t.Error("Indexer not added to list of indexers")
	}
}

func TestStatus(t *testing.T) {
	ind1 := MockTVIndexer{}
	ind2 := MockMovieIndexer{}

	indexersCollection = []Indexer{ind1, ind2}

	mods, err := Status()
	if len(mods) != 2 {
		t.Errorf("Expected to have 2 indexers modules status, got %d instead", len(mods))
	}
	if err == nil {
		t.Error("Expected to have aggregated error for indexer status")
	}
}

func TestGetTorrentForEpisode(t *testing.T) {
	ind1 := MockTVIndexer{}
	ind2 := MockTVIndexer{}
	ind3 := MockMovieIndexer{}
	ind4 := MockMovieIndexer{}

	indexersCollection = []Indexer{ind1, ind2, ind3, ind4}

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

func TestGetTorrentForMovie(t *testing.T) {
	ind1 := MockTVIndexer{}
	ind2 := MockTVIndexer{}
	ind3 := MockMovieIndexer{}
	ind4 := MockMovieIndexer{}

	indexersCollection = []Indexer{ind1, ind2, ind3, ind4}

	torrentList, _ := GetTorrentForMovie("Test movie")
	if len(torrentList) != 6 {
		t.Errorf("Expected 6 torrents, got %d instead\n", len(torrentList))
	}

	if torrentList[0].Seeders != 3 {
		t.Error("Torrent list is not sorted by seeders")
	}

	torrentList, err := GetTorrentForMovie("")
	if err == nil {
		t.Error("Expected to have zero results and an error, go no error instead: ")
	}
	if len(torrentList) != 0 {
		t.Errorf("Expected to have no results, got %d results instead\n", len(torrentList))
	}
}
