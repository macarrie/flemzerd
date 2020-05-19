package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/macarrie/flemzerd/configuration"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/stats"
)

var srv *http.Server
var router *gin.Engine
var serverStarted bool

var identityKey = "id"

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	UserName string
}

func init() {
	serverStarted = false
}

func initRouter() {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "flemzerd",
		Key:         []byte("flemzerd_key"),
		Timeout:     336 * time.Hour, // 2 Weeks
		MaxRefresh:  0,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				log.Error("Error", err)
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if userID == configuration.Config.Interface.Auth.Username && password == configuration.Config.Interface.Auth.Password {
				return &User{
					UserName: userID,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	router = gin.Default()
	router.Use(static.Serve("/static", static.LocalFile("/var/lib/flemzerd/server/ui/static", true)))
	router.Use(static.Serve("/cache", static.LocalFile(CACHE_BASE_PATH, true)))

	router.NoRoute(func(c *gin.Context) {
		c.File("/var/lib/flemzerd/server/ui/index.html")
	})

	routes := router.Group("/")
	{
		routes.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/dashboard")
		})

		// Manifest and service worker are required for PWA
		routes.GET("/manifest.json", func(c *gin.Context) {
			c.File("/var/lib/flemzerd/server/ui/manifest.json")
		})

		routes.GET("/favicon.ico", func(c *gin.Context) {
			c.File("/var/lib/flemzerd/server/ui/favicon.ico")
		})

		// Register js files route for pwa (precache-manifest.js and others)
		webuiRoot := "/var/lib/flemzerd/server/ui/"
		err := filepath.Walk(webuiRoot, func(path string, info os.FileInfo, err error) error {
			extension := filepath.Ext(path)
			if extension == ".js" || extension == ".png" {
				routes.GET(filepath.Base(path), func(c *gin.Context) {
					c.File(path)
				})
			}
			return nil
		})
		if err != nil {
			log.Error(err)
		}
	}

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.GET("/refresh_token", authMiddleware.RefreshHandler)
			auth.POST("/login", authMiddleware.LoginHandler)
		}

		configRoute := v1.Group("/config")
		{
			configRoute.GET("/", authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
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
		actionsRoute.Use(authMiddleware.MiddlewareFunc())
		{
			actionsRoute.POST("/poll", actionsCheckNow)
			actionsRoute.POST("/refresh_all_metadata", resetAllMetadata)
		}

		tvshowsRoute := v1.Group("/tvshows")
		tvshowsRoute.Use(authMiddleware.MiddlewareFunc())
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
			tvshowsRoute.DELETE("/details/:id", deleteShow)
			tvshowsRoute.POST("/restore/:id", restoreShow)
			tvshowsRoute.GET("/episodes/:id", getEpisodeDetails)
			tvshowsRoute.POST("/episodes/:id/download", downloadEpisode)
			tvshowsRoute.DELETE("/episodes/:id", deleteEpisode)
			tvshowsRoute.DELETE("/episodes/:id/download", abortEpisodeDownload)
			tvshowsRoute.POST("/episodes/:id/download/skip_torrent", skipEpisodeTorrentDownload)
			tvshowsRoute.PUT("/episodes/:id/download_state", changeEpisodeDownloadedState)
			tvshowsRoute.POST("/details/:id/refresh_metadata", refreshShowMetadata)
		}

		moviesRoute := v1.Group("/movies")
		moviesRoute.Use(authMiddleware.MiddlewareFunc())
		{
			moviesRoute.GET("/tracked", getTrackedMovies)
			moviesRoute.GET("/downloading", getDownloadingMovies)
			moviesRoute.GET("/removed", getRemovedMovies)
			moviesRoute.GET("/downloaded", getDownloadedMovies)
			moviesRoute.GET("/details/:id", getMovieDetails)
			moviesRoute.DELETE("/details/:id", deleteMovie)
			moviesRoute.POST("/details/:id/download", downloadMovie)
			moviesRoute.DELETE("/details/:id/download", abortMovieDownload)
			moviesRoute.POST("/details/:id/download/skip_torrent", skipMovieTorrentDownload)
			moviesRoute.PUT("/details/:id", updateMovie)
			moviesRoute.PUT("/details/:id/download_state", changeMovieDownloadedState)
			moviesRoute.PUT("/details/:id/custom_title", changeMovieCustomTitle)
			moviesRoute.PUT("/details/:id/use_default_title", useMovieDefaultTitle)
			moviesRoute.POST("/restore/:id", restoreMovie)
			moviesRoute.POST("/details/:id/refresh_metadata", refreshMovieMetadata)
		}

		modules := v1.Group("/modules")
		modules.Use(authMiddleware.MiddlewareFunc())
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

			scanner := modules.Group("/scan")
			{
				movies_scan_group := scanner.Group("/movies")
				{
					movies_scan_group.GET("/", scanMovieLibrary)
					movies_scan_group.POST("/import", importScannedMovie)
				}
				show_scan_group := scanner.Group("/tvshows")
				{
					show_scan_group.GET("/", scanShowLibrary)
					show_scan_group.POST("/import", importScannedShow)
				}
			}
		}

		notificationsRoute := v1.Group("/notifications")
		notificationsRoute.Use(authMiddleware.MiddlewareFunc())
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

func Start(port int) {
	//gin.SetMode(gin.DebugMode)
	gin.SetMode(gin.ReleaseMode)

	initRouter()

	var srv *http.Server
	if configuration.Config.Interface.UseSSL {
		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
		srv = &http.Server{
			Addr:         fmt.Sprintf("0.0.0.0:%d", port),
			Handler:      router,
			TLSConfig:    tlsConfig,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
	} else {
		srv = &http.Server{
			Addr:    fmt.Sprintf("0.0.0.0:%d", port),
			Handler: router,
		}
	}

	log.WithFields(log.Fields{
		"port": port,
	}).Info("Starting HTTP server")

	serverStarted = true

	var err error
	if configuration.Config.Interface.UseSSL {
		err = srv.ListenAndServeTLS(configuration.Config.Interface.SSLCert, configuration.Config.Interface.SSLServerKey)
	} else {
		err = srv.ListenAndServe()
	}

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
