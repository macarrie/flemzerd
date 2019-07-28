package db

import (
	"testing"

	downloadable "github.com/macarrie/flemzerd/downloadable"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

func init() {
	DbFatal = func(code int) {
		log.Error("Database fatal error: ", code)
	}

	DbPath = "/tmp/flemzerd.db"
	Load()
}

func TestDBReadOnly(t *testing.T) {
	DbPath = "../testdata/not_writable.db"
	Load()

	DbPath = "/tmp/flemzerd.db"
	Load()
}

func TestGetTrackedItems(t *testing.T) {
	ResetDb()
	t1 := TvShow{
		Title: "t1",
	}
	t2 := TvShow{
		Title: "t2",
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
		Title: "e1",
		DownloadingItem: DownloadingItem{
			Downloading: true,
		},
	}
	e2 := Episode{
		Title: "e2",
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
		Title: "e1",
		DownloadingItem: DownloadingItem{
			Downloaded: true,
		},
	}
	e2 := Episode{
		Title: "e2",
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

func TestGetRemovedItems(t *testing.T) {
	ResetDb()

	m := Movie{
		Title:         "removed_movie",
		OriginalTitle: "removed_movie",
	}
	show := TvShow{
		Title:         "removed_show",
		OriginalTitle: "removed_show",
	}

	Client.Create(&m)
	Client.Delete(&m)

	Client.Create(&show)
	Client.Delete(&show)

	movieList, _ := GetRemovedMovies()
	showList, _ := GetRemovedTvShows()

	if len(movieList) != 1 {
		t.Errorf("Expected 1 item in removed movies list, got %d instead", len(movieList))
	}
	if len(showList) != 1 {
		t.Errorf("Expected 1 item in removed movies list, got %d instead", len(movieList))
	}
}

func TestGetNotifications(t *testing.T) {
	ResetDb()

	n1 := Notification{
		Title: "notification_1",
		Read:  false,
	}
	n2 := Notification{
		Title: "notification_2",
		Read:  true,
	}

	Client.Create(&n1)
	Client.Create(&n2)

	notificationsList, _ := GetNotifications()
	readNotificationsList, _ := GetReadNotifications()
	unreadNotificationsList, _ := GetUnreadNotifications()

	if len(notificationsList) != 2 {
		t.Errorf("Expected 2 items in notifications list, got %d instead", len(notificationsList))
	}
	if len(readNotificationsList) != 1 {
		t.Errorf("Expected 1 item in read notifications list, got %d instead", len(readNotificationsList))
	}
	if len(unreadNotificationsList) != 1 {
		t.Errorf("Expected 1 item in unread notifications list, got %d instead", len(unreadNotificationsList))
	}
}

func TestSaveTraktAndTelegramInfos(t *testing.T) {
	SaveTraktToken("test")
	SaveTelegramChatID(1234)
}

func TestSaveDownloadable(t *testing.T) {
	ResetDb()

	movie := Movie{
		Title:         "test_movie_save_downloadable",
		OriginalTitle: "test_movie_save_downloadable",
	}

	movie.DownloadingItem.Downloading = true
	movie.DownloadingItem = DownloadingItem{
		TorrentList: []Torrent{
			Torrent{
				Name: "test_torrent_1",
			},
			Torrent{
				Name: "test_torrent_2",
			},
			Torrent{
				Name:   "test_torrent_3",
				Failed: true,
			},
			Torrent{
				Name:   "test_torrent_4",
				Failed: true,
			},
		},
	}
	var movieDownloadable downloadable.Downloadable = &movie
	SaveDownloadable(&movieDownloadable)

	var movieFromDB Movie
	Client.Where(Movie{Title: "test_movie_save_downloadable"}).First(&movieFromDB)
	if movieFromDB.DownloadingItem.CurrentTorrent().ID == 0 || len(movieFromDB.DownloadingItem.TorrentList) != 4 || len(movieFromDB.DownloadingItem.FailedTorrents()) != 2 {
		t.Error("Expected movie downloading item to be saved during SaveDownloadable")
	}

	episode := Episode{
		Title: "test_episode_save_downloadable",
	}

	test_torrent := Torrent{
		Name: "test_torrent",
	}
	Client.Save(&test_torrent)
	episode.DownloadingItem.Downloading = true
	episode.DownloadingItem = DownloadingItem{
		TorrentList: []Torrent{
			Torrent{
				Name: "test_torrent_1",
			},
			Torrent{
				Name: "test_torrent_2",
			},
			Torrent{
				Name:   "test_torrent_3",
				Failed: true,
			},
			Torrent{
				Name:   "test_torrent_4",
				Failed: true,
			},
		},
	}
	var episodeDownloadable downloadable.Downloadable = &episode
	SaveDownloadable(&episodeDownloadable)

	var episodeFromDB Episode
	Client.Where(Episode{Title: "test_episode_save_downloadable"}).First(&episodeFromDB)
	if episodeFromDB.DownloadingItem.CurrentTorrent().ID == 0 || len(episodeFromDB.DownloadingItem.TorrentList) != 4 || len(episodeFromDB.DownloadingItem.FailedTorrents()) != 2 {
		t.Error("Expected episode downloading item to be saved during SaveDownloadable")
	}
}
