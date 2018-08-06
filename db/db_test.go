package db

import (
	"strconv"
	"testing"

	. "github.com/macarrie/flemzerd/objects"
)

func init() {
	DbPath = "/tmp/flemzerd.db"
	Load()
}

func TestTorrentIsFailed(t *testing.T) {
	torrent := Torrent{
		TorrentId: strconv.Itoa(TORRENT_STOPPED),
		Name:      "Test torrent",
		Link:      "test.torrent",
	}
	movie := Movie{
		Title:         "movie",
		OriginalTitle: "movie",
		DownloadingItem: DownloadingItem{
			Downloading: true,
		},
	}

	if TorrentHasFailed(movie.DownloadingItem, torrent) {
		t.Error("Expected torrent not to be in failed torrents")
	}

	movie.DownloadingItem.FailedTorrents = []Torrent{torrent}

	if !TorrentHasFailed(movie.DownloadingItem, torrent) {
		t.Error("Expected torrent to be in failed torrents")
	}
}

func TestGetTrackedItems(t *testing.T) {
	ResetDb()
	t1 := TvShow{
		Name: "t1",
	}
	t2 := TvShow{
		Name: "t2",
	}
	m1 := Movie{
		Title:         "m1",
		OriginalTitle: "m1",
		DownloadingItem: DownloadingItem{
			Downloading: true,
		},
	}
	m2 := Movie{
		Title:         "m2",
		OriginalTitle: "m2",
	}
	m3 := Movie{
		Title:         "m3",
		OriginalTitle: "m3",
		DownloadingItem: DownloadingItem{
			Downloading: true,
		},
	}

	Client.Create(&t1)
	Client.Create(&t2)
	Client.Create(&m1)
	Client.Create(&m2)
	Client.Create(&m3)

	showList, _ := GetTrackedTvShows()
	movieList, _ := GetTrackedMovies()

	if len(showList) != 2 {
		t.Errorf("Expected 2 items in tracked shows list, got %d instead", len(showList))
	}
	if len(movieList) != 1 {
		t.Errorf("Expected 1 items in tracked movies list, got %d instead", len(movieList))
	}
}

func TestGetDownloadingItems(t *testing.T) {
	ResetDb()
	e1 := Episode{
		Name: "e1",
		DownloadingItem: DownloadingItem{
			Downloading: true,
		},
	}
	e2 := Episode{
		Name: "e2",
		DownloadingItem: DownloadingItem{
			Downloading: true,
		},
	}
	m1 := Movie{
		Title:         "m1",
		OriginalTitle: "m1",
		DownloadingItem: DownloadingItem{
			Downloading: true,
		},
	}
	m2 := Movie{
		Title:         "m2",
		OriginalTitle: "m2",
		DownloadingItem: DownloadingItem{
			Downloading: true,
		},
	}

	Client.Create(&e1)
	Client.Create(&e2)
	Client.Create(&m1)
	Client.Create(&m2)

	episodeList, _ := GetDownloadingEpisodes()
	movieList, _ := GetDownloadingMovies()

	if len(episodeList) != 2 {
		t.Errorf("Expected 2 items in downloading episodes list, got %d instead", len(episodeList))
	}
	if len(movieList) != 2 {
		t.Errorf("Expected 2 items in downloading episodes list, got %d instead", len(movieList))
	}
}

func TestGetDownloadedItems(t *testing.T) {
	ResetDb()
	e1 := Episode{
		Name: "e1",
		DownloadingItem: DownloadingItem{
			Downloaded: true,
		},
	}
	e2 := Episode{
		Name: "e2",
		DownloadingItem: DownloadingItem{
			Downloaded: true,
		},
	}
	m1 := Movie{
		Title:         "m1",
		OriginalTitle: "m1",
		DownloadingItem: DownloadingItem{
			Downloaded: true,
		},
	}
	m2 := Movie{
		Title:         "m2",
		OriginalTitle: "m2",
		DownloadingItem: DownloadingItem{
			Downloaded: true,
		},
	}

	Client.Create(&e1)
	Client.Create(&e2)
	Client.Create(&m1)
	Client.Create(&m2)

	episodeList, _ := GetDownloadedEpisodes()
	movieList, _ := GetDownloadedMovies()

	if len(episodeList) != 2 {
		t.Errorf("Expected 2 items in downloaded episodes list, got %d instead", len(episodeList))
	}
	if len(movieList) != 2 {
		t.Errorf("Expected 2 items in downloaded episodes list, got %d instead", len(movieList))
	}
}

func TestTraktToken(t *testing.T) {
	token := "test"

	SaveTraktToken(token)

	savedToken := LoadTraktToken()
	if savedToken != token {
		t.Errorf("Expected token from db to be %s, got %s instead", token, savedToken)
	}
}
