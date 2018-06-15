package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
)

var srv *http.Server
var router *gin.Engine
var serverStarted bool

func init() {
	serverStarted = false
}

func initRouter() {
	router = gin.Default()
	router.Use(static.Serve("/", static.LocalFile("/var/lib/flemzerd/server/ui", true)))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/config", func(c *gin.Context) {
			c.JSON(http.StatusOK, configuration.Config)
		})

		tvshowsRoute := v1.Group("/tvshows")
		{
			tvshowsRoute.GET("/tracked", getTrackedShows)
			tvshowsRoute.GET("/downloading", getDownloadingEpisodes)
			tvshowsRoute.GET("/removed", getRemovedShows)
			tvshowsRoute.GET("/downloaded", getDownloadedEpisodes)
			tvshowsRoute.GET("/details/:id", getShowDetails)
			tvshowsRoute.GET("/details/:id/seasons/:season_nb", getSeasonDetails)
			tvshowsRoute.DELETE("/details/:id", deleteShow)
			tvshowsRoute.POST("/restore/:id", restoreShow)
			tvshowsRoute.GET("/episodes/:id", getEpisodeDetails)
			tvshowsRoute.POST("/episodes/:id/download", downloadEpisode)
			tvshowsRoute.DELETE("/episodes/:id", deleteEpisode)
			tvshowsRoute.DELETE("/episodes/:id/download", abortEpisodeDownload)
			tvshowsRoute.PUT("/episodes/:id", changeEpisodeDownloadedState)
		}

		moviesRoute := v1.Group("/movies")
		{
			moviesRoute.GET("/tracked", getTrackedMovies)
			moviesRoute.GET("/downloading", getDownloadingMovies)
			moviesRoute.GET("/removed", getRemovedMovies)
			moviesRoute.GET("/downloaded", getDownloadedMovies)
			moviesRoute.GET("/details/:id", getMovieDetails)
			moviesRoute.DELETE("/details/:id", deleteMovie)
			moviesRoute.DELETE("/details/:id/download", abortMovieDownload)
			moviesRoute.PUT("/details/:id", changeMovieDownloadedState)
			moviesRoute.POST("/restore/:id", restoreMovie)
		}

		modules := v1.Group("/modules")
		{
			modules.GET("/status", getModulesStatus)

			providers := modules.Group("/providers")
			{
				providers.GET("/status", getProvidersStatus)
			}

			indexers := modules.Group("/indexers")
			{
				indexers.GET("/status", getIndexersStatus)
			}

			notifiers := modules.Group("/notifiers")
			{
				notifiers.GET("/status", getNotifiersStatus)
			}

			downloaders := modules.Group("/downloaders")
			{
				downloaders.GET("/status", getDownloadersStatus)
			}

			mediacenters := modules.Group("/mediacenters")
			{
				mediacenters.GET("/status", getMediacentersStatus)
			}

			watchlists := modules.Group("/watchlists")
			{
				watchlists.GET("/status", getWatchlistsStatus)
				watchlists.POST("/refresh", refreshWatchlists)

				traktRoutes := watchlists.Group("/trakt")
				{
					traktRoutes.GET("/auth", performTraktAuth)
					traktRoutes.GET("/auth_errors", getTraktAuthErrors)
					traktRoutes.GET("/token", getTraktToken)
					traktRoutes.GET("/devicecode", getTraktDeviceCode)
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

	serverStarted = true

	err := srv.ListenAndServe()
	if err != nil {
		log.Error(fmt.Sprintf("listen: %s\n", err))
		serverStarted = false
		return
	}
}

func Stop() {
	if !serverStarted {
		log.Debug("No webserver to stop")
		return
	}

	log.Info("Stopping HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cancel()
	srv.Shutdown(ctx)
	serverStarted = false
}
