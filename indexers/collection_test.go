package indexer

import (
	"fmt"
	"testing"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	mock "github.com/macarrie/flemzerd/mocks"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/vidocq"
)

func init() {
	// go test makes a cd into package directory when testing. We must go up by one level to load our testdata
	configuration.UseFile("../testdata/test_config.toml")
	err := configuration.Load()
	if err != nil {
		fmt.Print("Could not load test configuration file: ", err)
	}
}

func TestAddIndexer(t *testing.T) {
	AddIndexer(mock.TVIndexer{})
	if len(indexersCollection) != 1 {
		t.Error("Indexer not added to list of indexers")
	}
}

func TestStatus(t *testing.T) {
	ind1 := mock.TVIndexer{}
	ind2 := mock.ErrorMovieIndexer{}

	indexersCollection = []Indexer{ind1, ind2}

	mods, err := Status()
	if len(mods) != 2 {
		t.Errorf("Expected to have 2 indexers modules status, got %d instead", len(mods))
	}
	if err == nil {
		t.Error("Expected to have aggregated error for indexer status")
	}

	ind3 := mock.TVIndexer{}
	indexersCollection = []Indexer{ind3}
	mods, err = Status()
	if err != nil {
		t.Error("Expected to have no errors when getting status from ok indexers")
	}
}

func TestReset(t *testing.T) {
	ind := mock.MovieIndexer{}
	AddIndexer(ind)
	Reset()

	if len(indexersCollection) != 0 {
		t.Error("Expected indexer collection to be empty after reset")
	}
}

func TestGetTorrentForEpisode(t *testing.T) {
	ind1 := mock.TVIndexer{}
	ind2 := mock.TVIndexer{}
	ind3 := mock.MovieIndexer{}
	ind4 := mock.MovieIndexer{}
	configuration.Config.System.PreferredMediaQuality = "720p,1080p"
	configuration.Config.System.ExcludedReleaseTypes = "cam,screener,telesync,telecine"

	indexersCollection = []Indexer{ind1, ind2, ind3, ind4}

	torrentList, _ := GetTorrentForEpisode("Test show", 1, 1)
	if len(torrentList) != 6 {
		t.Errorf("Expected 6 torrents, got %d instead\n", len(torrentList))
		return
	}

	if torrentList[0].Seeders != 3 {
		t.Error("Torrent list is not sorted by seeders")
	}

	torrentList, err := GetTorrentForEpisode("Test show", 1, 0)
	if err != nil {
		t.Error("Expected to have zero results and no error, got an error instead: ")
	}
	if len(torrentList) != 0 {
		t.Errorf("Expected to have no results, got %d results instead\n", len(torrentList))
	}

	torrentList, _ = GetTorrentForEpisode("Test show", 0, 1)
	if len(torrentList) > 0 {
		t.Error("Expected to have no torrents when getting torrents for episode")
	}

	//Get torrent when vidocq is not available
	vidocq.LocalVidocqAvailable = false
	torrentList, _ = GetTorrentForEpisode("Test show", 1, 1)
	if len(torrentList) != 14 {
		t.Errorf("Expected 14 torrents, got %d instead\n", len(torrentList))
		return
	}

	//Get torrents with filters disabled
	configuration.Config.System.PreferredMediaQuality = ""
	configuration.Config.System.ExcludedReleaseTypes = ""
	vidocq.LocalVidocqAvailable = true
	torrentList, _ = GetTorrentForEpisode("Test show", 1, 1)
	if len(torrentList) != 8 {
		t.Errorf("Expected 6 torrents, got %d instead\n", len(torrentList))
		return
	}
}

func TestGetTorrentForMovie(t *testing.T) {
	ind1 := mock.TVIndexer{}
	ind2 := mock.TVIndexer{}
	ind3 := mock.MovieIndexer{}
	ind4 := mock.MovieIndexer{}
	configuration.Config.System.PreferredMediaQuality = "720p,1080p"
	configuration.Config.System.ExcludedReleaseTypes = "cam,screener,telesync,telecine"
	vidocq.LocalVidocqAvailable = true

	movieDate := time.Date(2018, time.January, 10, 13, 0, 0, 0, time.UTC)

	indexersCollection = []Indexer{ind1, ind2, ind3, ind4}

	torrentList, _ := GetTorrentForMovie(Movie{
		Title:         "Test movie",
		OriginalTitle: "Test movie",
		Date:          movieDate,
	})
	if len(torrentList) != 6 {
		t.Errorf("Expected 6 torrents, got %d instead\n", len(torrentList))
	}

	if torrentList[0].Seeders != 3 {
		t.Error("Torrent list is not sorted by seeders")
	}

	torrentList, err := GetTorrentForMovie(Movie{
		Date: movieDate,
	})
	if err != nil {
		t.Error("Expected to have zero results and no error, got an error instead: ")
	}
	if len(torrentList) != 0 {
		t.Errorf("Expected to have no results, got %d results instead\n", len(torrentList))
	}

	torrentList, _ = GetTorrentForMovie(Movie{
		Title:         "error",
		OriginalTitle: "error",
	})
	if len(torrentList) > 0 {
		t.Error("Expected to have no torrents when getting torrents for movie")
	}

	//Get torrent when vidocq is not available
	vidocq.LocalVidocqAvailable = false
	torrentList, _ = GetTorrentForMovie(Movie{
		Title:         "Test movie",
		OriginalTitle: "Test movie",
		Date:          movieDate,
	})
	if len(torrentList) != 14 {
		t.Errorf("Expected 14 torrents when vidocq is not available, got %d instead\n", len(torrentList))
	}

	//Get torrents with no filters
	configuration.Config.System.PreferredMediaQuality = ""
	configuration.Config.System.ExcludedReleaseTypes = ""
	vidocq.LocalVidocqAvailable = true
	torrentList, _ = GetTorrentForMovie(Movie{
		Title:         "Test movie",
		OriginalTitle: "Test movie",
		Date:          movieDate,
	})
	if len(torrentList) != 12 {
		t.Errorf("Expected 12 torrents, got %d instead\n", len(torrentList))
	}

}

func TestGetSpecificIndexer(t *testing.T) {
	p := mock.TVIndexer{}
	indexersCollection = []Indexer{p}

	if _, err := GetIndexer("Unknown"); err == nil {
		t.Error("Expected to have error when getting unknown indexer, got none")
	}

	if _, err := GetIndexer("TVIndexer"); err != nil {
		t.Errorf("Got error while retrieving known indexer: %s", err.Error())
	}
}
