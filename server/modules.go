package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/flemzerd/downloaders"
	"github.com/macarrie/flemzerd/indexers"
	"github.com/macarrie/flemzerd/mediacenters"
	"github.com/macarrie/flemzerd/notifiers"
	"github.com/macarrie/flemzerd/providers"
	"github.com/macarrie/flemzerd/watchlists"
)

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
