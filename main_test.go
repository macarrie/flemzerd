package main

import (
	"testing"

	log "github.com/macarrie/flemzerd/logging"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	"github.com/macarrie/flemzerd/downloaders"
	"github.com/macarrie/flemzerd/indexers"
	"github.com/macarrie/flemzerd/mediacenters"
	"github.com/macarrie/flemzerd/notifiers"
	"github.com/macarrie/flemzerd/providers"
	"github.com/macarrie/flemzerd/watchlists"
)

func init() {
	log.Setup(true)

	db.DbPath = "/tmp/flemzerd.db"
	db.Load()
	db.ResetDb()

	// go test makes a cd into package directory when testing. We must go up by one level to load our testdata
	configuration.UseFile("testdata/test_config.toml")
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

	initConfiguration()

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

	initConfiguration()

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
