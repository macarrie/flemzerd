package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	watchlist "github.com/macarrie/flemzerd/watchlists"
	"github.com/macarrie/flemzerd/watchlists/impl/trakt"
)

func performTraktAuth(c *gin.Context) {
	w, err := watchlist.GetWatchlist("trakt")
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	t := w.(*trakt.TraktWatchlist)
	if err := t.IsAuthenticated(); err == nil {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	go t.Auth()
	c.JSON(http.StatusOK, gin.H{})
	return
}

func getTraktAuthErrors(c *gin.Context) {
	w, err := watchlist.GetWatchlist("trakt")
	t := w.(*trakt.TraktWatchlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, t.GetAuthErrors())
}

func getTraktToken(c *gin.Context) {
	w, err := watchlist.GetWatchlist("trakt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	t := w.(*trakt.TraktWatchlist)
	c.JSON(http.StatusOK, t.Token)
}

func getTraktDeviceCode(c *gin.Context) {
	w, err := watchlist.GetWatchlist("trakt")
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	t := w.(*trakt.TraktWatchlist)
	c.JSON(http.StatusOK, t.DeviceCode)
	return
}
