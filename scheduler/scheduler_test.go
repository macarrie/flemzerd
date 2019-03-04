package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/macarrie/flemzerd/db"

	"github.com/macarrie/flemzerd/configuration"

	downloader "github.com/macarrie/flemzerd/downloaders"
	indexer "github.com/macarrie/flemzerd/indexers"
	mediacenter "github.com/macarrie/flemzerd/mediacenters"
	notifier "github.com/macarrie/flemzerd/notifiers"
	provider "github.com/macarrie/flemzerd/providers"
	watchlist "github.com/macarrie/flemzerd/watchlists"

	log "github.com/macarrie/flemzerd/logging"
	mock "github.com/macarrie/flemzerd/mocks"
	. "github.com/macarrie/flemzerd/objects"
)

func init() {
	log.Setup(true)

	db.DbPath = "/tmp/flemzerd.db"
	db.Load()
	db.ResetDb()

	// go test makes a cd into package directory when testing. We must go up by one level to load our testdata
	configuration.UseFile("../testdata/test_config.toml")
}

func TestInitConfiguration(t *testing.T) {
	providerList := []string{
		"TMDB",
		"TVDB",
	}
	notifierList := []string{
		"Desktop notifier",
		"Pushbullet",
		"telegram",
	}
	indexerList := []string{
		"Indexer 1",
		"Indexer 2",
	}
	downloaderList := []string{
		"transmission",
	}
	watchlistList := []string{
		"manual",
		"trakt",
	}
	mediacenterList := []string{
		"kodi",
	}

	initConfiguration(true)

	for _, not := range notifierList {
		if _, err := notifier.GetNotifier(not); err != nil {
			t.Errorf("Expected to have '%s' notifier loaded", not)
		}
	}
	for _, prov := range providerList {
		if _, err := provider.GetProvider(prov); err != nil {
			t.Errorf("Expected to have '%s' provider loaded", prov)
		}
	}
	for _, ind := range indexerList {
		if _, err := indexer.GetIndexer(ind); err != nil {
			t.Errorf("Expected to have '%s' indexer loaded", ind)
		}
	}
	for _, dl := range downloaderList {
		if _, err := downloader.GetDownloader(dl); err != nil {
			t.Errorf("Expected to have '%s' downloader loaded", dl)
		}
	}
	for _, wl := range watchlistList {
		if _, err := watchlist.GetWatchlist(wl); err != nil {
			t.Errorf("Expected to have '%s' watchlist loaded", wl)
		}
	}
	for _, md := range mediacenterList {
		if _, err := mediacenter.GetMediaCenter(md); err != nil {
			t.Errorf("Expected to have '%s' mediacenter loaded", md)
		}
	}

	notifier.Reset()
	provider.Reset()
	indexer.Reset()
	downloader.Reset()
	watchlist.Reset()
	mediacenter.Reset()

	Reload(true)

	for _, not := range notifierList {
		if _, err := notifier.GetNotifier(not); err != nil {
			t.Errorf("Expected to have '%s' notifier loaded", not)
		}
	}
	for _, prov := range providerList {
		if _, err := provider.GetProvider(prov); err != nil {
			t.Errorf("Expected to have '%s' provider loaded", prov)
		}
	}
	for _, ind := range indexerList {
		if _, err := indexer.GetIndexer(ind); err != nil {
			t.Errorf("Expected to have '%s' indexer loaded", ind)
		}
	}
	for _, dl := range downloaderList {
		if _, err := downloader.GetDownloader(dl); err != nil {
			t.Errorf("Expected to have '%s' downloader loaded", dl)
		}
	}
	for _, wl := range watchlistList {
		if _, err := watchlist.GetWatchlist(wl); err != nil {
			t.Errorf("Expected to have '%s' watchlist loaded", wl)
		}
	}
	for _, md := range mediacenterList {
		if _, err := mediacenter.GetMediaCenter(md); err != nil {
			t.Errorf("Expected to have '%s' mediacenter loaded", md)
		}
	}
}

func TestDownloadMedia(t *testing.T) {
	notifier.Reset()
	provider.Reset()
	indexer.Reset()
	downloader.Reset()
	watchlist.Reset()
	mediacenter.Reset()

	provider.AddProvider(mock.TVProvider{})
	provider.AddProvider(mock.MovieProvider{})
	indexer.AddIndexer(mock.ErrorTVIndexer{})
	indexer.AddIndexer(mock.ErrorMovieIndexer{})
	downloader.AddDownloader(mock.Downloader{})
	notifier.AddNotifier(mock.Notifier{})
	watchlist.AddWatchlist(mock.Watchlist{})
	mediacenter.AddMediaCenter(mock.MediaCenter{})

	episode := Episode{
		Title:  "test_episode",
		Number: 0,
		Season: 1,
		TvShow: TvShow{
			Title:         "test_show",
			OriginalTitle: "test_show",
		},
	}
	db.Client.Save(&episode)

	// Test case: No torrents found
	Download(&episode, false)

	db.Client.Find(&episode, episode.ID)
	if !episode.DownloadingItem.TorrentsNotFound {
		t.Error("Expected download to fail because no torrent can be found")
	}

	// Test case: Error while getting torrents
	episode.Number = 1
	episode.Season = 0
	episode.DownloadingItem = DownloadingItem{}
	db.Client.Save(&episode)

	Download(&episode, false)

	db.Client.Find(&episode, episode.ID)
	if !episode.DownloadingItem.Pending || episode.DownloadingItem.Downloading {
		t.Error("Expected download to fail because no torrent can be requested")
	}

	// Test case: Episode already downloading
	episode.DownloadingItem = DownloadingItem{}
	episode.DownloadingItem.Downloading = true
	db.Client.Save(&episode)

	Download(&episode, false)

	db.Client.Find(&episode, episode.ID)
	if episode.DownloadingItem.Pending || !episode.DownloadingItem.Downloading {
		t.Error("Expected download to do nothing because episode is already downloading")
	}

	// Test case: Episode already downloading
	episode.DownloadingItem = DownloadingItem{}
	episode.DownloadingItem.Downloading = false
	episode.DownloadingItem.Downloaded = true
	db.Client.Save(&episode)

	Download(&episode, false)

	db.Client.Find(&episode, episode.ID)
	if episode.DownloadingItem.Pending || episode.DownloadingItem.Downloading || !episode.DownloadingItem.Downloaded {
		t.Error("Expected download to do nothing because episode is already downloaded")
	}

	movie := Movie{
		Title:         "",
		OriginalTitle: "",
	}
	db.Client.Save(&movie)

	// Test case: No torrents found
	Download(&movie, false)

	db.Client.Find(&movie, movie.ID)
	if !movie.DownloadingItem.TorrentsNotFound {
		t.Error("Expected download to fail because no torrent can be found")
	}

	// Test case: Error while getting torrents
	movie.OriginalTitle = "error"
	db.Client.Save(&movie)

	Download(&movie, false)

	db.Client.Find(&movie, movie.ID)
	if !movie.DownloadingItem.Pending || movie.DownloadingItem.Downloading {
		t.Error("Expected download to fail because no torrent can be requested")
	}

	// Test case: movie already downloading
	movie.DownloadingItem = DownloadingItem{}
	movie.DownloadingItem.Downloading = true
	db.Client.Save(&movie)

	Download(&movie, false)

	db.Client.Find(&movie, movie.ID)
	if movie.DownloadingItem.Pending || !movie.DownloadingItem.Downloading {
		t.Error("Expected download to do nothing because movie is already downloading")
	}

	// Test case: movie already downloading
	movie.DownloadingItem = DownloadingItem{}
	movie.DownloadingItem.Downloading = false
	movie.DownloadingItem.Downloaded = true
	db.Client.Save(&movie)

	Download(&movie, false)

	db.Client.Find(&movie, movie.ID)
	if movie.DownloadingItem.Pending || movie.DownloadingItem.Downloading || !movie.DownloadingItem.Downloaded {
		t.Error("Expected download to do nothing because movie is already downloaded")
	}
}

func TestPoll(t *testing.T) {
	// Test with error modules

	notifier.Reset()
	provider.Reset()
	indexer.Reset()
	downloader.Reset()
	watchlist.Reset()
	mediacenter.Reset()

	provider.AddProvider(mock.ErrorProvider{})
	indexer.AddIndexer(mock.ErrorTVIndexer{})
	indexer.AddIndexer(mock.ErrorMovieIndexer{})
	downloader.AddDownloader(mock.ErrorDownloader{})
	notifier.AddNotifier(mock.ErrorNotifier{})
	watchlist.AddWatchlist(mock.ErrorWatchlist{})
	mediacenter.AddMediaCenter(mock.ErrorMediaCenter{})

	recoveryDone := false
	poll(&recoveryDone)

	// Test with correct modules
	notifier.Reset()
	provider.Reset()
	indexer.Reset()
	downloader.Reset()
	watchlist.Reset()
	mediacenter.Reset()

	provider.AddProvider(mock.TVProvider{})
	provider.AddProvider(mock.MovieProvider{})
	indexer.AddIndexer(mock.TVIndexer{})
	indexer.AddIndexer(mock.MovieIndexer{})
	downloader.AddDownloader(mock.Downloader{})
	notifier.AddNotifier(mock.Notifier{})
	watchlist.AddWatchlist(mock.Watchlist{})
	mediacenter.AddMediaCenter(mock.MediaCenter{})

	poll(&recoveryDone)

	for id, ctxStorage := range downloader.EpisodeDownloadRoutines {
		ctx, timeout := context.WithTimeout(ctxStorage.Context, 10*time.Second)
		defer timeout()

		select {
		case <-ctx.Done():
			err := ctx.Err()
			switch err {
			case context.Canceled:
				var episode Episode
				db.Client.Find(&episode, id)
				if !episode.DownloadingItem.Downloaded {
					t.Error("Expected episode to be marked as downloaded when download finished")
				}
				return
			case context.DeadlineExceeded:
				t.Error("Expected episode to finish download (timeout)")
			}
		}
	}

	for id, ctxStorage := range downloader.MovieDownloadRoutines {
		ctx, timeout := context.WithTimeout(ctxStorage.Context, 10*time.Second)
		defer timeout()

		select {
		case <-ctx.Done():
			err := ctx.Err()
			switch err {
			case context.Canceled:
				var movie Movie
				db.Client.Find(&movie, id)
				if !movie.DownloadingItem.Downloaded {
					t.Error("Expected movie to be marked as downloaded when download finished")
				}
				return
			case context.DeadlineExceeded:
				t.Error("Expected movie to finish download (timeout)")
			}
		}

	}
}

func TestRunAndStop(t *testing.T) {
	notifier.Reset()
	provider.Reset()
	indexer.Reset()
	downloader.Reset()
	watchlist.Reset()
	mediacenter.Reset()

	provider.AddProvider(mock.TVProvider{})
	provider.AddProvider(mock.MovieProvider{})
	indexer.AddIndexer(mock.TVIndexer{})
	indexer.AddIndexer(mock.MovieIndexer{})
	downloader.AddDownloader(mock.Downloader{})
	notifier.AddNotifier(mock.Notifier{})
	watchlist.AddWatchlist(mock.Watchlist{})
	mediacenter.AddMediaCenter(mock.MediaCenter{})

	configuration.Config.System.CheckInterval = 1

	go Run(true)
	time.Sleep(10 * time.Second)
	Stop()

	db.Load()
}

func TestRecovery(t *testing.T) {
	db.ResetDb()
	downloader.EpisodeDownloadRoutines = make(map[uint](downloader.ContextStorage))
	downloader.MovieDownloadRoutines = make(map[uint](downloader.ContextStorage))

	notifier.Reset()
	provider.Reset()
	indexer.Reset()
	downloader.Reset()
	watchlist.Reset()
	mediacenter.Reset()

	provider.AddProvider(mock.TVProvider{})
	provider.AddProvider(mock.MovieProvider{})
	indexer.AddIndexer(mock.TVIndexer{})
	indexer.AddIndexer(mock.MovieIndexer{})
	downloader.AddDownloader(mock.Downloader{})
	notifier.AddNotifier(mock.Notifier{})
	watchlist.AddWatchlist(mock.Watchlist{})
	mediacenter.AddMediaCenter(mock.MediaCenter{})

	provider.GetTVShowsInfoFromConfig()
	provider.GetMoviesInfoFromConfig()

	for _, show := range provider.TVShows {
		recentEpisodes, _ := provider.FindRecentlyAiredEpisodesForShow(show)

		for _, recentEpisode := range recentEpisodes {
			reqEpisode := Episode{}
			req := db.Client.Where(Episode{
				Title:  recentEpisode.Title,
				Season: recentEpisode.Season,
				Number: recentEpisode.Number,
			}).Find(&reqEpisode)
			if req.RecordNotFound() {
				recentEpisode.TvShow = show
				db.Client.Create(&recentEpisode)
			} else {
				recentEpisode = reqEpisode
			}

			recentEpisode.DownloadingItem.Downloading = true
			db.Client.Save(&recentEpisode)
		}
	}
	for _, movie := range provider.Movies {
		movie.DownloadingItem.Downloading = true
		db.Client.Save(&movie)
	}

	RecoverDownloadingItems()

	for id, ctxStorage := range downloader.EpisodeDownloadRoutines {
		log.Error("DEBUG EPISODE ID: ", id)
		ctx, timeout := context.WithTimeout(ctxStorage.Context, 10*time.Second)
		defer timeout()

		select {
		case <-ctx.Done():
			err := ctx.Err()
			switch err {
			case context.Canceled:
				log.Error("DEBUG CANCELLED")
				var episode Episode
				db.Client.Find(&episode, id)
				if !episode.DownloadingItem.Downloaded {
					t.Error("Expected episode to be marked as downloaded when download finished")
				}
				return
			case context.DeadlineExceeded:
				log.Error("DEBUG TIMEOUT")
				t.Error("Expected episode to finish download (timeout)")
			}
		}
	}
	for id, ctxStorage := range downloader.MovieDownloadRoutines {
		ctx, timeout := context.WithTimeout(ctxStorage.Context, 10*time.Second)
		defer timeout()

		select {
		case <-ctx.Done():
			err := ctx.Err()
			switch err {
			case context.Canceled:
				var movie Movie
				db.Client.Find(&movie, id)
				if !movie.DownloadingItem.Downloaded {
					t.Error("Expected movie to be marked as downloaded when download finished")
				}
				return
			case context.DeadlineExceeded:
				t.Error("Expected movie to finish download (timeout)")
			}
		}

	}
}
