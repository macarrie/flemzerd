package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/flemzerd/logging"

	"github.com/macarrie/flemzerd/db"

	. "github.com/macarrie/flemzerd/objects"
)

func getTrackedMovies(c *gin.Context) {
	movies, err := db.GetTrackedMovies()
	if err != nil {
		log.Error("Error while getting downloading movies from db: ", err)
	}
	c.JSON(http.StatusOK, movies)
}

func getDownloadingMovies(c *gin.Context) {
	movies, err := db.GetDownloadingMovies()
	if err != nil {
		log.Error("Error while getting downloading movies from db: ", err)
	}
	c.JSON(http.StatusOK, movies)
}

func getRemovedMovies(c *gin.Context) {
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
}

func getDownloadedMovies(c *gin.Context) {
	movies, err := db.GetDownloadedMovies()
	if err != nil {
		log.Error("Error while getting downloaded movies from db: ", err)
	}
	c.JSON(http.StatusOK, movies)
}

func getMovieDetails(c *gin.Context) {
	id := c.Param("id")
	var movie Movie
	req := db.Client.Unscoped().Find(&movie, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, movie)
}

func deleteMovie(c *gin.Context) {
	id := c.Param("id")
	var movie Movie
	req := db.Client.Delete(&movie, id)
	if err := req.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func changeMovieDownloadedState(c *gin.Context) {
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
}

func restoreMovie(c *gin.Context) {
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
}
