package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/watchlists/impl/trakt"

	"github.com/macarrie/flemzerd/configuration"
	downloader "github.com/macarrie/flemzerd/downloaders"
	indexer "github.com/macarrie/flemzerd/indexers"
	notifier "github.com/macarrie/flemzerd/notifiers"
	provider "github.com/macarrie/flemzerd/providers"
	watchlist "github.com/macarrie/flemzerd/watchlists"

	. "github.com/macarrie/flemzerd/objects"
)

var srv *http.Server
var router *gin.Engine

func initRouter() {
	router = gin.Default()

	router.Static("/static", "/var/lib/flemzerd/server/ui/")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui")
	})

	router.LoadHTMLFiles("/var/lib/flemzerd/server/ui/index.html")
	router.GET("/ui", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	v1 := router.Group("/api/v1")
	{
		v1.GET("/config", func(c *gin.Context) {
			c.JSON(http.StatusOK, configuration.Config)
		})

		v1.GET("/tvshows", func(c *gin.Context) {
			c.JSON(http.StatusOK, provider.TVShows)
		})

		v1.GET("/movies", func(c *gin.Context) {
			c.JSON(http.StatusOK, provider.Movies)
		})

		modules := v1.Group("/modules")
		{
			modules.GET("/status", func(c *gin.Context) {
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
			})

			providers := modules.Group("/providers")
			{
				providers.GET("/status", func(c *gin.Context) {
					mods, _ := provider.Status()
					c.JSON(http.StatusOK, mods)
				})
			}

			indexers := modules.Group("/indexers")
			{
				indexers.GET("/status", func(c *gin.Context) {
					mods, _ := indexer.Status()
					c.JSON(http.StatusOK, mods)
				})
			}

			notifiers := modules.Group("/notifiers")
			{
				notifiers.GET("/status", func(c *gin.Context) {
					mods, _ := notifier.Status()
					c.JSON(http.StatusOK, mods)
				})
			}

			downloaders := modules.Group("/downloaders")
			{
				downloaders.GET("/status", func(c *gin.Context) {
					mods, _ := downloader.Status()
					c.JSON(http.StatusOK, mods)
				})
			}

			watchlists := modules.Group("/watchlists")
			{
				watchlists.GET("/status", func(c *gin.Context) {
					mods, _ := watchlist.Status()
					c.JSON(http.StatusOK, mods)
				})

				traktRoutes := watchlists.Group("/trakt")
				{

					traktRoutes.GET("/auth", func(c *gin.Context) {
						w, err := watchlist.GetWatchlist("Trakt")
						t := w.(*trakt.TraktWatchlist)
						if err == nil {
							if err := t.IsAuthenticated(); err == nil {
								c.JSON(http.StatusNoContent, gin.H{})
								return
							}

							go t.Auth()
							c.JSON(http.StatusOK, gin.H{})
							return
						}
						c.JSON(http.StatusNotFound, err)
					})

					traktRoutes.GET("/auth_errors", func(c *gin.Context) {
						w, err := watchlist.GetWatchlist("Trakt")
						t := w.(*trakt.TraktWatchlist)
						if err != nil {
							c.JSON(http.StatusInternalServerError, err)
							return
						}

						c.JSON(http.StatusOK, t.GetAuthErrors())
					})

					traktRoutes.GET("/token", func(c *gin.Context) {
						w, err := watchlist.GetWatchlist("Trakt")
						t := w.(*trakt.TraktWatchlist)
						if err != nil {
							c.JSON(http.StatusInternalServerError, err)
							return
						}

						c.JSON(http.StatusOK, t.Token)
					})

					traktRoutes.GET("/devicecode", func(c *gin.Context) {
						w, err := watchlist.GetWatchlist("Trakt")
						t := w.(*trakt.TraktWatchlist)
						if err == nil {
							c.JSON(http.StatusOK, t.DeviceCode)
							return
						}
						c.JSON(http.StatusNotFound, err)
					})
				}
			}
		}

	}

}

func Start(port int, debug bool) {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	initRouter()

	srv = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: router,
	}

	log.WithFields(log.Fields{
		"port": port,
	}).Info("Starting HTTP server")
	err := srv.ListenAndServe()
	if err != nil {
		log.Debug(fmt.Sprintf("listen: %s\n", err))
	}
}

func Stop() {
	log.Info("Stopping HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cancel()
	srv.Shutdown(ctx)
}