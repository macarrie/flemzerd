package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	provider "github.com/macarrie/flemzerd/providers"
	"github.com/macarrie/flemzerd/scanner"

	. "github.com/macarrie/flemzerd/objects"
)

func scanMovieLibrary(c *gin.Context) {
	log.Error("SCANNING MOVIE LIBRARY")
	movies, err := scanner.ScanMovies()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not scan movie library correctly")
	}

	c.JSON(http.StatusOK, movies)
}

func scanShowLibrary(c *gin.Context) {
	log.Error("SCANNING SHOW LIBRARY")
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
	_, err := provider.FindMovie(idsFromDb)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func importScannedShow(c *gin.Context) {
	var showInfoFromRequest MediaInfo
	c.BindJSON(&showInfoFromRequest)

	idsFromDb := MediaIds{}
	req := db.Client.Where("title = ?", showInfoFromRequest.Title).Find(&idsFromDb)
	if req.RecordNotFound() {
		idsFromDb = MediaIds{
			Title: showInfoFromRequest.Title,
		}
	}
	_, err := provider.FindShow(idsFromDb)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
