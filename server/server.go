package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/stats"
)

var srv *http.Server
var router *gin.Engine
var serverStarted bool

func init() {
	serverStarted = false
}

func initRouter() {
	router = gin.Default()
	router.Use(static.Serve("/static", static.LocalFile("/var/lib/flemzerd/server/ui/static", true)))

	router.NoRoute(func(c *gin.Context) {
		c.File("/var/lib/flemzerd/server/ui/index.html")
	})

	routes := router.Group("/")
	{
		routes.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/dashboard")
		})
	}

	v1 := router.Group("/api/v1")
	{
		configRoute := v1.Group("/config")
		{
			configRoute.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, configuration.Config)
			})

			configRoute.GET("/check", func(c *gin.Context) {
				if err := configuration.Check(); err != nil {
					if multierr, ok := err.(*multierror.Error); ok {
						c.JSON(http.StatusOK, multierr.Errors)
						return
					}
					c.JSON(http.StatusOK, gin.H{})
					return
				}
				c.JSON(http.StatusOK, gin.H{})
				return
			})
		}

		v1.GET("/stats", stats.Handler())

		actionsRoute := v1.Group("/actions")
		{
			actionsRoute.POST("/poll", actionsCheckNow)
		}

		tvshowsRoute := v1.Group("/tvshows")
		{
			tvshowsRoute.GET("/tracked", getTrackedShows)
			tvshowsRoute.GET("/downloading", getDownloadingEpisodes)
			tvshowsRoute.GET("/removed", getRemovedShows)
			tvshowsRoute.GET("/downloaded", getDownloadedEpisodes)
			tvshowsRoute.GET("/details/:id", getShowDetails)
			tvshowsRoute.PUT("/details/:id", updateShow)
			tvshowsRoute.PUT("/details/:id/custom_title", changeTvshowCustomTitle)
			tvshowsRoute.PUT("/details/:id/use_default_title", useTvshowDefaultTitle)
			tvshowsRoute.PUT("/details/:id/change_anime_state", changeTvshowAnimeState)
			tvshowsRoute.GET("/details/:id/seasons/:season_nb", getSeasonDetails)
			tvshowsRoute.POST("/details/:id/seasons/:season_nb/download", downloadSeason)
			tvshowsRoute.PUT("/details/:id/seasons/:season_nb/download_state", changeSeasonDownloadedState)
			tvshowsRoute.DELETE("/details/:id", deleteShow)
			tvshowsRoute.POST("/restore/:id", restoreShow)
			tvshowsRoute.GET("/episodes/:id", getEpisodeDetails)
			tvshowsRoute.POST("/episodes/:id/download", downloadEpisode)
			tvshowsRoute.DELETE("/episodes/:id", deleteEpisode)
			tvshowsRoute.DELETE("/episodes/:id/download", abortEpisodeDownload)
			tvshowsRoute.PUT("/episodes/:id/download_state", changeEpisodeDownloadedState)
		}

		moviesRoute := v1.Group("/movies")
		{
			moviesRoute.GET("/tracked", getTrackedMovies)
			moviesRoute.GET("/downloading", getDownloadingMovies)
			moviesRoute.GET("/removed", getRemovedMovies)
			moviesRoute.GET("/downloaded", getDownloadedMovies)
			moviesRoute.GET("/details/:id", getMovieDetails)
			moviesRoute.DELETE("/details/:id", deleteMovie)
			moviesRoute.POST("/details/:id/download", downloadMovie)
			moviesRoute.DELETE("/details/:id/download", abortMovieDownload)
			moviesRoute.PUT("/details/:id", updateMovie)
			moviesRoute.PUT("/details/:id/download_state", changeMovieDownloadedState)
			moviesRoute.PUT("/details/:id/custom_title", changeMovieCustomTitle)
			moviesRoute.PUT("/details/:id/use_default_title", useMovieDefaultTitle)
			moviesRoute.POST("/restore/:id", restoreMovie)
		}

		modules := v1.Group("/modules")
		{
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

				telegramRoutes := notifiers.Group("/telegram")
				{
					telegramRoutes.GET("/auth", performTelegramAuth)
					telegramRoutes.GET("/chatid", getTelegramChatID)
					telegramRoutes.GET("/auth_code", getTelegramAuthCode)
				}
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

		notificationsRoute := v1.Group("/notifications")
		{
			notificationsRoute.GET("/all", getNotifications)
			notificationsRoute.DELETE("/all", deleteNotifications)
			notificationsRoute.GET("/read", getReadNotifications)
			notificationsRoute.POST("/read", markAllNotificationsAsRead)
			notificationsRoute.POST("/read/:id", changeNotificationReadState)
			notificationsRoute.GET("/unread", getUnreadNotifications)
			notificationsRoute.POST("/unread/:id", changeNotificationReadState)
			notificationsRoute.GET("/render/list", renderNotificationsList)
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
