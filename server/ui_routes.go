package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"

	downloader "github.com/macarrie/flemzerd/downloaders"
	indexer "github.com/macarrie/flemzerd/indexers"
	mediacenter "github.com/macarrie/flemzerd/mediacenters"
	notifier "github.com/macarrie/flemzerd/notifiers"
	provider "github.com/macarrie/flemzerd/providers"
	watchlist "github.com/macarrie/flemzerd/watchlists"

	log "github.com/macarrie/flemzerd/logging"

	. "github.com/macarrie/flemzerd/objects"
)

func ui_dashboard(c *gin.Context) {
	trackedMovies, err := db.GetTrackedMovies()
	if err != nil {
		log.Error("Error while getting downloading movies from db: ", err)
	}

	downloadingMovies, err := db.GetDownloadingMovies()
	if err != nil {
		log.Error("Error while getting downloading movies from db: ", err)
	}

	downloadedMovies, err := db.GetDownloadedMovies()
	if err != nil {
		log.Error("Error while getting downloaded movies from db: ", err)
	}

	trackedShows, err := db.GetTrackedTvShows()
	if err != nil {
		log.Error("Error while gettings tracked shows from db: ", err)
	}
	downloadingEpisodes, err := db.GetDownloadingEpisodes()
	if err != nil {
		log.Error("Error while getting downloading episodes from db: ", err)
	}
	downloadedEpisodes, err := db.GetDownloadedEpisodes()
	if err != nil {
		log.Error("Error while getting downloaded episodes from db: ", err)
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"trackedMoviesNb":       len(trackedMovies),
		"downloadingMoviesNb":   len(downloadingMovies),
		"downloadedMoviesNb":    len(downloadedMovies),
		"trackedShowsNb":        len(trackedShows),
		"downloadingEpisodesNb": len(downloadingEpisodes),
		"downloadedEpisodesNb":  len(downloadedEpisodes),
	})
}

func ui_tvshows(c *gin.Context) {
	trackedShows, err := db.GetTrackedTvShows()
	if err != nil {
		log.Error("Error while gettings tracked shows from db: ", err)
	}
	downloadingEpisodes, err := db.GetDownloadingEpisodes()
	if err != nil {
		log.Error("Error while getting downloading episodes from db: ", err)
	}

	var tvShows []TvShow
	var removedShows []TvShow

	if err := db.Client.Unscoped().Order("status").Order("title").Find(&tvShows).Error; err != nil {
		log.Error("Error while getting removed shows from db: ", err)
	}

	for _, show := range tvShows {
		if show.DeletedAt != nil {
			removedShows = append(removedShows, show)
		}
	}

	c.HTML(http.StatusOK, "tvshows", gin.H{
		"trackedShows":        trackedShows,
		"downloadingEpisodes": downloadingEpisodes,
		//"downloadedEpisodes":  downloadedEpisodes,
		"removedShows": removedShows,
	})
}

func ui_tvshow(c *gin.Context) {
	id := c.Param("id")
	var show TvShow
	req := db.Client.Unscoped().Find(&show, id)
	if req.RecordNotFound() {
		c.HTML(http.StatusNotFound, "404", gin.H{})
		return
	}

	var seasonsDetails []SeasonDetails
	for _, season := range show.Seasons {
		epList, err := provider.GetSeasonEpisodeList(show, season.SeasonNumber)
		if err != nil {
			log.WithFields(log.Fields{
				"show":   show.OriginalTitle,
				"season": season.SeasonNumber,
				"error":  err,
			}).Warning("Encountered error when querying season details")
			c.JSON(http.StatusInternalServerError, gin.H{})
			continue
		}
		seasonsDetails = append(seasonsDetails, SeasonDetails{
			Info:        season,
			EpisodeList: epList,
		})
	}

	c.HTML(http.StatusOK, "tvshow_details", gin.H{
		"item":           show,
		"seasonsdetails": seasonsDetails,
		"type":           "tvshow",
	})
}

func ui_episode(c *gin.Context) {
	id := c.Param("id")
	var ep Episode
	req := db.Client.Find(&ep, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.HTML(http.StatusOK, "episode_details", gin.H{
		"item": ep,
		"type": "episode",
	})
}

func ui_movies(c *gin.Context) {
	trackedMovies, err := db.GetTrackedMovies()
	if err != nil {
		log.Error("Error while getting downloading movies from db: ", err)
	}

	downloadingMovies, err := db.GetDownloadingMovies()
	if err != nil {
		log.Error("Error while getting downloading movies from db: ", err)
	}

	downloadedMovies, err := db.GetDownloadedMovies()
	if err != nil {
		log.Error("Error while getting downloaded movies from db: ", err)
	}

	var removedMoviesQuery []Movie
	var removedMovies []Movie

	if err := db.Client.Unscoped().Order("created_at DESC").Find(&removedMoviesQuery).Error; err != nil {
		log.Error("Error while getting removed movies from db: ", err)
	}

	for _, m := range removedMoviesQuery {
		if m.DeletedAt != nil {
			removedMovies = append(removedMovies, m)
		}
	}

	c.HTML(http.StatusOK, "movies", gin.H{
		"trackedMovies":     trackedMovies,
		"downloadingMovies": downloadingMovies,
		"downloadedMovies":  downloadedMovies,
		"removedMovies":     removedMovies,
		"config":            configuration.Config,
	})
}

func ui_movie(c *gin.Context) {
	id := c.Param("id")
	var movie Movie
	req := db.Client.Unscoped().Find(&movie, id)
	if req.RecordNotFound() {
		c.HTML(http.StatusNotFound, "404", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "movie_details", gin.H{
		"item": movie,
		"type": "movie",
	})
}

func ui_status(c *gin.Context) {
	providers, _ := provider.Status()
	notifiers, _ := notifier.Status()
	indexers, _ := indexer.Status()
	downloaders, _ := downloader.Status()
	mediacenters, _ := mediacenter.Status()
	watchlists, _ := watchlist.Status()

	c.HTML(http.StatusOK, "status", gin.H{
		"providers":    providers,
		"notifiers":    notifiers,
		"indexers":     indexers,
		"downloaders":  downloaders,
		"mediacenters": mediacenters,
		"watchlists":   watchlists,
	})
}

func ui_settings(c *gin.Context) {
	c.HTML(http.StatusOK, "settings", gin.H{
		"config": configuration.Config,
	})
}

func ui_notifications(c *gin.Context) {
	var notifs []Notification
	db.Client.Order("created_at DESC").Find(&notifs)

	c.HTML(http.StatusOK, "notifications", gin.H{
		"notifications": notifs,
	})
}
