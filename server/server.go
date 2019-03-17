package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/macarrie/flemzerd/configuration"
	notifier_helper "github.com/macarrie/flemzerd/helpers/notifiers"

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

	router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:      "/var/lib/flemzerd/server/ui/templates",
		Extension: ".tpl",
		Master:    "layout/base",
		Partials: []string{
			"partials/movie_miniature",
			"partials/tvshow_miniature",
			"partials/module_status",
			"partials/downloading_item",
			"partials/notifications_list",
			"partials/notification",
			"partials/media_details_action_bar",
			"partials/media_ids",
		},
		Funcs: template.FuncMap{
			"isInFuture": func(date time.Time) bool {
				return date.After(time.Now())
			},
			"getNotificationType": notifier_helper.GetNotificationType,
		},
		DisableCache: true,
	})

	router.Use(static.Serve("/static", static.LocalFile("/var/lib/flemzerd/server/ui/static", true)))

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404", gin.H{})
	})

	routes := router.Group("/")
	{
		routes.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/dashboard")
		})

		routes.GET("/dashboard", ui_dashboard)
		routes.GET("/tvshows", ui_tvshows)
		routes.GET("/tvshows/:id", ui_tvshow)
		routes.GET("/episodes/:id", ui_episode)
		routes.GET("/movies", ui_movies)
		routes.GET("/movies/:id", ui_movie)
		routes.GET("/status", ui_status)
		routes.GET("/settings", ui_settings)
		routes.GET("/notifications", ui_notifications)
	}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/config", func(c *gin.Context) {
			c.JSON(http.StatusOK, configuration.Config)
		})

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
