package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	downloader "github.com/macarrie/flemzerd/downloaders"
	indexer "github.com/macarrie/flemzerd/indexers"
	log "github.com/macarrie/flemzerd/logging"
	notifier "github.com/macarrie/flemzerd/notifiers"
	provider "github.com/macarrie/flemzerd/providers"

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
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, epList)
}

func deleteShow(c *gin.Context) {
	id := c.Param("id")
	var show TvShow
	req := db.Client.Delete(&show, id)
	if err := req.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
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
	c.JSON(http.StatusOK, ep)
}

func downloadEpisode(c *gin.Context) {
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
		c.JSON(http.StatusNotFound, gin.H{"error": "No torrents found"})
		return
	}
	notifier.NotifyEpisodeDownloadStart(&ep)

	go downloader.DownloadEpisode(ep.TvShow, ep, toDownload)

	c.JSON(http.StatusOK, gin.H{})
	return
}

func deleteEpisode(c *gin.Context) {
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
	c.AbortWithStatus(http.StatusNoContent)
}

func changeEpisodeDownloadedState(c *gin.Context) {
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
}
