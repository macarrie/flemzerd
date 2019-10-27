package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	provider "github.com/macarrie/flemzerd/providers"
	"github.com/macarrie/flemzerd/scanner"

	. "github.com/macarrie/flemzerd/objects"
)

func scanMovieLibrary(c *gin.Context) {
	movies, err := scanner.ScanMovies()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not scan movie library correctly")
	}

	c.JSON(http.StatusOK, movies)
}

func scanShowLibrary(c *gin.Context) {
	shows, err := scanner.ScanShows()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not scan show library correctly")
	}

	c.JSON(http.StatusOK, shows)
}

func importScannedMovie(c *gin.Context) {
	var movieInfoFromRequest MediaInfo
	c.BindJSON(&movieInfoFromRequest)

	idsFromDb := MediaIds{}
	req := db.Client.Where("title = ?", movieInfoFromRequest.Title).Find(&idsFromDb)
	if req.RecordNotFound() {
		idsFromDb = MediaIds{
			Title: movieInfoFromRequest.Title,
		}
	}
	movie, err := provider.FindMovie(idsFromDb)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	torrentNamePath := strings.Split(movieInfoFromRequest.Raw, "/")
	torrentName := torrentNamePath[len(torrentNamePath)-1]
	movie.DownloadingItem = DownloadingItem{
		Downloaded: true,
		TorrentList: []Torrent{
			Torrent{
				Name: torrentName,
			},
		},
	}
	db.Client.Save(&movie)
	log.WithFields(log.Fields{
		"movie": movie.GetTitle(),
	}).Info("Imported movie from scan")

	c.JSON(http.StatusOK, gin.H{})
}

func importScannedShow(c *gin.Context) {
	var showInfoFromRequest struct {
		Title string
		Data  MediaInfoSeasons
	}
	c.BindJSON(&showInfoFromRequest)

	idsFromDb := MediaIds{}
	req := db.Client.Where("title = ?", showInfoFromRequest.Title).Find(&idsFromDb)
	if req.RecordNotFound() {
		idsFromDb = MediaIds{
			Title: showInfoFromRequest.Title,
		}
	}
	show, err := provider.FindShow(idsFromDb)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	season_from_provider := make(map[int][]Episode)
	for nb, season := range showInfoFromRequest.Data {
		season_from_provider[nb], err = provider.GetSeasonEpisodeList(show, nb)
		if err != nil {
			log.WithFields(log.Fields{
				"show":   show.GetTitle(),
				"season": nb,
				"error":  err,
			}).Error("Could not get details for TV Show season when trying to import episodes from scan. Skipping all episodes from this season")
			continue
		}

		for episode_nb := range season {
			for episode := range season_from_provider[nb] {
				if episode_nb == season_from_provider[nb][episode].Number {
					torrentNamePath := strings.Split(season[episode_nb].Raw, "/")
					torrentName := torrentNamePath[len(torrentNamePath)-1]
					episode_from_provider := season_from_provider[nb][episode]
					episode_from_provider.DownloadingItem = DownloadingItem{
						Downloaded: true,
						TorrentList: []Torrent{
							Torrent{
								Name: torrentName,
							},
						},
					}
					db.Client.Save(&episode_from_provider)
					log.WithFields(log.Fields{
						"show":    show.GetTitle(),
						"season":  nb,
						"episode": episode_nb,
					}).Info("Imported episode from scan")
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}
