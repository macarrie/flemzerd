package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/flemzerd/logging"
	mediacenter "github.com/macarrie/flemzerd/mediacenters"
	"github.com/macarrie/flemzerd/watchlists/impl/trakt"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	downloader "github.com/macarrie/flemzerd/downloaders"
	indexer "github.com/macarrie/flemzerd/indexers"
	notifier "github.com/macarrie/flemzerd/notifiers"
	provider "github.com/macarrie/flemzerd/providers"
	watchlist "github.com/macarrie/flemzerd/watchlists"

	. "github.com/macarrie/flemzerd/objects"
)

var srv *http.Server
var router *gin.Engine
var serverStarted bool

func init() {
	serverStarted = false
}

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

		tvshowsRoute := v1.Group("/tvshows")
		{
			tvshowsRoute.GET("/tracked", func(c *gin.Context) {
				shows, err := db.GetTrackedTvShows()
				if err != nil {
					log.Error("Error while gettings tracked shows from db: ", err)
				}
				c.JSON(http.StatusOK, shows)
			})
			tvshowsRoute.GET("/downloading", func(c *gin.Context) {
				episodes, err := db.GetDownloadingEpisodes()
				if err != nil {
					log.Error("Error while getting downloading episodes from db: ", err)
				}
				c.JSON(http.StatusOK, episodes)
			})
			tvshowsRoute.GET("/removed", func(c *gin.Context) {
				var tvShows []TvShow
				var retList []TvShow

				if err := db.Client.Unscoped().Order("status").Order("name").Find(&tvShows).Error; err != nil {
					log.Error("Error while getting removed shows from db: ", err)
				}

				for _, show := range tvShows {
					if show.DeletedAt != nil {
						retList = append(retList, show)
					}
				}
				c.JSON(http.StatusOK, retList)
			})
			tvshowsRoute.GET("/downloaded", func(c *gin.Context) {
				episodes, err := db.GetDownloadedEpisodes()
				if err != nil {
					log.Error("Error while getting downloaded episodes from db: ", err)
				}
				c.JSON(http.StatusOK, episodes)
			})
			tvshowsRoute.GET("/details/:id", func(c *gin.Context) {
				id := c.Param("id")
				var show TvShow
				req := db.Client.Unscoped().Find(&show, id)
				if req.RecordNotFound() {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}
				c.JSON(http.StatusOK, show)
			})
			tvshowsRoute.GET("/details/:id/seasons/:season_nb", func(c *gin.Context) {
				id := c.Param("id")
				seasonNumber := c.Param("season_nb")
				var show TvShow
				req := db.Client.First(&show, id)
				if req.RecordNotFound() {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}

				seasonNb, err := strconv.Atoi(seasonNumber)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": "Bad season number"})
					return
				}

				epList, err := provider.GetSeasonEpisodeList(show, seasonNb)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{})
					return
				}

				c.JSON(http.StatusOK, epList)
			})
			tvshowsRoute.DELETE("/details/:id", func(c *gin.Context) {
				id := c.Param("id")
				var show TvShow
				req := db.Client.Delete(&show, id)
				if err := req.Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}
				c.JSON(http.StatusNoContent, nil)
			})
			tvshowsRoute.POST("/restore/:id", func(c *gin.Context) {
				id := c.Param("id")
				var show TvShow
				req := db.Client.Unscoped().Find(&show, id)
				if req.RecordNotFound() {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}

				show.DeletedAt = nil
				db.Client.Unscoped().Save(&show)

				c.JSON(http.StatusOK, gin.H{})
			})
			tvshowsRoute.GET("/episodes/:id", func(c *gin.Context) {
				id := c.Param("id")
				var ep Episode
				req := db.Client.Find(&ep, id)
				if req.RecordNotFound() {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}
				c.JSON(http.StatusOK, ep)
			})
			tvshowsRoute.POST("/episodes/:id/download", func(c *gin.Context) {
				id := c.Param("id")

				var ep Episode
				req := db.Client.Find(&ep, id)
				if req.RecordNotFound() {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}

				if ep.DownloadingItem.Downloaded || ep.DownloadingItem.Downloading || ep.DownloadingItem.Pending {
					c.JSON(http.StatusNotModified, gin.H{})
					return
				}

				ep.DownloadingItem.Pending = true
				db.Client.Save(&ep)

				go func() {
					log.WithFields(log.Fields{
						"id":      id,
						"show":    ep.TvShow.Name,
						"episode": ep.Name,
						"season":  ep.Season,
						"number":  ep.Number,
					}).Info("Launching individual episode download")

					torrentList, err := indexer.GetTorrentForEpisode(ep.TvShow.Name, ep.Season, ep.Number)
					if err != nil {
						log.Warning(err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": err})
						return
					}
					log.Debug("Torrents found: ", len(torrentList))

					toDownload := downloader.FillEpisodeToDownloadTorrentList(&ep, torrentList)
					if len(toDownload) == 0 {
						downloader.MarkEpisodeFailedDownload(&ep.TvShow, &ep)
						c.JSON(http.StatusInternalServerError, gin.H{"error": "No torrents found"})
						return
					}
					notifier.NotifyEpisodeDownloadStart(&ep)

					downloader.DownloadEpisode(ep.TvShow, ep, toDownload)
				}()

				c.JSON(http.StatusOK, gin.H{})
				return
			})
			tvshowsRoute.DELETE("/episodes/:id", func(c *gin.Context) {
				id := c.Param("id")
				var ep Episode
				req := db.Client.Find(&ep, id)
				if req.RecordNotFound() {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}

				currentDownloadPath := ep.DownloadingItem.CurrentTorrent.DownloadDir
				downloader.RemoveTorrent(ep.DownloadingItem.CurrentTorrent)
				os.Remove(currentDownloadPath)

				db.Client.Unscoped().Delete(&ep)
				c.JSON(http.StatusNoContent, nil)
			})
			tvshowsRoute.PUT("/episodes/:id", func(c *gin.Context) {
				id := c.Param("id")
				var ep Episode
				req := db.Client.Find(&ep, id)
				if req.RecordNotFound() {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}

				var itemInfo DownloadingItem
				c.Bind(&itemInfo)

				ep.DownloadingItem.Downloaded = itemInfo.Downloaded
				db.Client.Save(&ep)

				c.JSON(http.StatusOK, ep)
			})
		}

		moviesRoute := v1.Group("/movies")
		{
			moviesRoute.GET("/tracked", func(c *gin.Context) {
				movies, err := db.GetTrackedMovies()
				if err != nil {
					log.Error("Error while getting downloading movies from db: ", err)
				}
				c.JSON(http.StatusOK, movies)
			})
			moviesRoute.GET("/downloading", func(c *gin.Context) {
				movies, err := db.GetDownloadingMovies()
				if err != nil {
					log.Error("Error while getting downloading movies from db: ", err)
				}
				c.JSON(http.StatusOK, movies)
			})
			moviesRoute.GET("/removed", func(c *gin.Context) {
				var movies []Movie
				var retList []Movie

				if err := db.Client.Unscoped().Order("created_at DESC").Find(&movies).Error; err != nil {
					log.Error("Error while getting removed movies from db: ", err)
				}

				for _, m := range movies {
					if m.DeletedAt != nil {
						retList = append(retList, m)
					}
				}
				c.JSON(http.StatusOK, retList)
			})
			moviesRoute.GET("/downloaded", func(c *gin.Context) {
				movies, err := db.GetDownloadedMovies()
				if err != nil {
					log.Error("Error while getting downloaded movies from db: ", err)
				}
				c.JSON(http.StatusOK, movies)
			})
			moviesRoute.GET("/details/:id", func(c *gin.Context) {
				id := c.Param("id")
				var movie Movie
				req := db.Client.Unscoped().Find(&movie, id)
				if req.RecordNotFound() {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}
				c.JSON(http.StatusOK, movie)
			})
			moviesRoute.DELETE("/details/:id", func(c *gin.Context) {
				id := c.Param("id")
				var movie Movie
				req := db.Client.Delete(&movie, id)
				if err := req.Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}
				c.JSON(http.StatusNoContent, nil)
			})
			moviesRoute.PUT("/details/:id", func(c *gin.Context) {
				id := c.Param("id")
				var movie Movie
				req := db.Client.Find(&movie, id)
				if req.RecordNotFound() {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}

				var itemInfo DownloadingItem
				c.Bind(&itemInfo)

				movie.DownloadingItem.Downloaded = itemInfo.Downloaded
				db.Client.Save(&movie)

				c.JSON(http.StatusOK, movie)
			})
			moviesRoute.POST("/restore/:id", func(c *gin.Context) {
				id := c.Param("id")
				var movie Movie
				req := db.Client.Unscoped().Find(&movie, id)
				if req.RecordNotFound() {
					c.JSON(http.StatusNotFound, gin.H{})
					return
				}

				movie.DeletedAt = nil
				db.Client.Unscoped().Save(&movie)

				c.JSON(http.StatusOK, gin.H{})
			})
		}

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

			mediacenters := modules.Group("/mediacenters")
			{
				mediacenters.GET("/status", func(c *gin.Context) {
					mods, _ := mediacenter.Status()
					c.JSON(http.StatusOK, mods)
				})
			}

			watchlists := modules.Group("/watchlists")
			{
				watchlists.GET("/status", func(c *gin.Context) {
					mods, _ := watchlist.Status()
					c.JSON(http.StatusOK, mods)
				})

				watchlists.POST("/refresh", func(c *gin.Context) {
					provider.GetTVShowsInfoFromConfig()
					provider.GetMoviesInfoFromConfig()
					c.JSON(http.StatusOK, gin.H{})
				})

				traktRoutes := watchlists.Group("/trakt")
				{

					traktRoutes.GET("/auth", func(c *gin.Context) {
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
					})

					traktRoutes.GET("/auth_errors", func(c *gin.Context) {
						w, err := watchlist.GetWatchlist("trakt")
						t := w.(*trakt.TraktWatchlist)
						if err != nil {
							c.JSON(http.StatusInternalServerError, err)
							return
						}

						c.JSON(http.StatusOK, t.GetAuthErrors())
					})

					traktRoutes.GET("/token", func(c *gin.Context) {
						w, err := watchlist.GetWatchlist("trakt")
						if err != nil {
							c.JSON(http.StatusInternalServerError, err)
							return
						}

						t := w.(*trakt.TraktWatchlist)
						c.JSON(http.StatusOK, t.Token)
					})

					traktRoutes.GET("/devicecode", func(c *gin.Context) {
						w, err := watchlist.GetWatchlist("trakt")
						if err != nil {
							c.JSON(http.StatusNotFound, err)
							return
						}

						t := w.(*trakt.TraktWatchlist)
						c.JSON(http.StatusOK, t.DeviceCode)
						return
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
