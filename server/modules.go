package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	downloader "github.com/macarrie/flemzerd/downloaders"
	indexer "github.com/macarrie/flemzerd/indexers"
	mediacenter "github.com/macarrie/flemzerd/mediacenters"
	notifier "github.com/macarrie/flemzerd/notifiers"
	provider "github.com/macarrie/flemzerd/providers"
	watchlist "github.com/macarrie/flemzerd/watchlists"

	. "github.com/macarrie/flemzerd/objects"
)

func getModulesStatus(c *gin.Context) {
	var status []Module

	providers, _ := provider.Status()
	notifiers, _ := notifier.Status()
	indexers, _ := indexer.Status()
	downloaders, _ := downloader.Status()

	status = append(status, providers...)
	status = append(status, notifiers...)
	status = append(status, indexers...)
	status = append(status, downloaders...)

	c.JSON(http.StatusOK, status)
}

func getProvidersStatus(c *gin.Context) {
	mods, _ := provider.Status()
	c.JSON(http.StatusOK, mods)
}

func getIndexersStatus(c *gin.Context) {
	mods, _ := indexer.Status()
	c.JSON(http.StatusOK, mods)
}

func getNotifiersStatus(c *gin.Context) {
	mods, _ := notifier.Status()
	c.JSON(http.StatusOK, mods)
}

func getDownloadersStatus(c *gin.Context) {
	mods, _ := downloader.Status()
	c.JSON(http.StatusOK, mods)
}

func getMediacentersStatus(c *gin.Context) {
	mods, _ := mediacenter.Status()
	c.JSON(http.StatusOK, mods)
}

func getWatchlistsStatus(c *gin.Context) {
	mods, _ := watchlist.Status()
	c.JSON(http.StatusOK, mods)
}

func refreshWatchlists(c *gin.Context) {
	provider.GetTVShowsInfoFromConfig()
	provider.GetMoviesInfoFromConfig()
	c.JSON(http.StatusOK, gin.H{})
}
