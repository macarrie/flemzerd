package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/macarrie/flemzerd/downloaders"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/providers"
	"github.com/macarrie/flemzerd/scheduler"
	"github.com/macarrie/flemzerd/stats"

	"github.com/macarrie/flemzerd/db"

	. "github.com/macarrie/flemzerd/objects"
)

func getTrackedShows(c *gin.Context) {
	shows, err := db.GetTrackedTvShows()
	if err != nil {
		log.Error("Error while gettings tracked shows from db: ", err)
	}
	c.JSON(http.StatusOK, shows)
}

func getDownloadingEpisodes(c *gin.Context) {
	episodes, err := db.GetDownloadingEpisodes()
	if err != nil {
		log.Error("Error while getting downloading episodes from db: ", err)
	}
	c.JSON(http.StatusOK, episodes)
}

func getRemovedShows(c *gin.Context) {
	tvshows, err := db.GetRemovedTvShows()
	if err != nil {
		log.Error("Error while getting removed tvshows from db: ", err)
	}
	c.JSON(http.StatusOK, tvshows)
}

func getDownloadedEpisodes(c *gin.Context) {
	episodes, err := db.GetDownloadedEpisodes()
	if err != nil {
		log.Error("Error while getting downloaded episodes from db: ", err)
	}
	c.JSON(http.StatusOK, episodes)
}

func getShowDetails(c *gin.Context) {
	id := c.Param("id")
	var show TvShow
	req := db.Client.Unscoped().Find(&show, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, show)
}

func updateShow(c *gin.Context) {
	id := c.Param("id")
	var show TvShow
	req := db.Client.Find(&show, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var showFromRequest TvShow
	c.Bind(&showFromRequest)

	show = showFromRequest
	db.Client.Save(&show)

	c.JSON(http.StatusOK, show)
}

func changeTvshowCustomTitle(c *gin.Context) {
	id := c.Param("id")
	var tvshow TvShow
	req := db.Client.Find(&tvshow, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var showFromRequest TvShow
	c.BindJSON(&showFromRequest)

	tvshow.CustomTitle = showFromRequest.CustomTitle
	db.Client.Save(&tvshow)

	c.JSON(http.StatusOK, tvshow)
}

func useTvshowDefaultTitle(c *gin.Context) {
	id := c.Param("id")
	var tvshow TvShow
	req := db.Client.Find(&tvshow, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var showFromRequest TvShow
	c.BindJSON(&showFromRequest)

	tvshow.UseDefaultTitle = showFromRequest.UseDefaultTitle
	db.Client.Save(&tvshow)

	c.JSON(http.StatusOK, tvshow)
}

func changeTvshowAnimeState(c *gin.Context) {
	id := c.Param("id")
	var tvshow TvShow
	req := db.Client.Find(&tvshow, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var showFromRequest TvShow
	c.BindJSON(&showFromRequest)

	tvshow.IsAnime = showFromRequest.IsAnime
	db.Client.Save(&tvshow)

	c.JSON(http.StatusOK, tvshow)
}

func getSeasonDetails(c *gin.Context) {
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
		log.WithFields(log.Fields{
			"show":   show.GetTitle(),
			"season": seasonNumber,
			"error":  err,
		}).Warning("Encountered error when querying season details from TMDB")
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	seasonInfo := TvSeason{}
	seasonInfoFound := false
	for _, season := range show.Seasons {
		if season.SeasonNumber == seasonNb {
			seasonInfoFound = true
			seasonInfo = season
		}
	}
	if !seasonInfoFound {
		c.JSON(http.StatusNotFound, gin.H{})
	}

	c.JSON(http.StatusOK, SeasonDetails{
		Info:        seasonInfo,
		EpisodeList: epList,
	})
}

func deleteShow(c *gin.Context) {
	id := c.Param("id")
	var show TvShow
	req := db.Client.Delete(&show, id)
	if err := req.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	stats.Stats.Shows.Tracked -= 1
	stats.Stats.Shows.Removed += 1
	c.AbortWithStatus(http.StatusNoContent)
}

func restoreShow(c *gin.Context) {
	id := c.Param("id")
	var show TvShow
	req := db.Client.Unscoped().Find(&show, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	show.DeletedAt = nil
	db.Client.Unscoped().Save(&show)

	stats.Stats.Shows.Tracked += 1
	stats.Stats.Shows.Removed -= 1
	c.JSON(http.StatusOK, gin.H{})
}

func getEpisodeDetails(c *gin.Context) {
	id := c.Param("id")
	var ep Episode
	req := db.Client.Find(&ep, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var show TvShow
	if ep.TvShow.ID == 0 && ep.TvShowID != 0 {
		// Episode tvshow has been deleted. Try to retrieve infos from db (show soft deleted)
		if req := db.Client.Unscoped().Find(&show, ep.TvShowID); req.RecordNotFound() {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		ep.TvShow = show
	}
	c.JSON(http.StatusOK, ep)
}

func downloadEpisode(c *gin.Context) {
	id := c.Param("id")

	var ep Episode
	req := db.Client.Find(&ep, id)
	if req.RecordNotFound() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if ep.DownloadingItem.Downloaded || ep.DownloadingItem.Downloading || ep.DownloadingItem.Pending {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}

	ep.GetLog().Info("Launching individual episode download")

	scheduler.Download(&ep, false)

	c.JSON(http.StatusOK, gin.H{})
	return
}

func deleteEpisode(c *gin.Context) {
	id := c.Param("id")
	var ep Episode
	req := db.Client.Delete(&ep, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func abortEpisodeDownload(c *gin.Context) {
	id := c.Param("id")
	var ep Episode
	req := db.Client.Find(&ep, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	downloader.AbortDownload(&ep)

	c.AbortWithStatus(http.StatusNoContent)
}

func changeEpisodeDownloadedState(c *gin.Context) {
	id := c.Param("id")
	var episode Episode
	req := db.Client.Find(&episode, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var downloadingItemFromRequest DownloadingItem
	c.BindJSON(&downloadingItemFromRequest)

	episode.DownloadingItem.Downloaded = downloadingItemFromRequest.Downloaded
	db.Client.Save(&episode)

	if episode.DownloadingItem.Downloaded {
		stats.Stats.Episodes.Downloaded += 1
	} else {
		stats.Stats.Episodes.Downloaded -= 1
	}

	c.JSON(http.StatusOK, episode)
}
