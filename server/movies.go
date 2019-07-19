package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/flemzerd/downloaders"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/scheduler"
	"github.com/macarrie/flemzerd/stats"

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
	movies, err := db.GetRemovedMovies()
	if err != nil {
		log.Error("Error while getting removed movies from db: ", err)
	}
	c.JSON(http.StatusOK, movies)
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
	req := db.Client.Unscoped().Find(&movie, id)
	if err := req.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	downloader.AbortDownload(&movie)
	req = db.Client.Delete(&movie, id)
	if err := req.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	stats.Stats.Movies.Tracked -= 1
	stats.Stats.Movies.Removed += 1
	c.AbortWithStatus(http.StatusNoContent)
}

func downloadMovie(c *gin.Context) {
	id := c.Param("id")

	var movie Movie
	req := db.Client.Find(&movie, id)
	if req.RecordNotFound() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if movie.DownloadingItem.Downloaded || movie.DownloadingItem.Downloading || movie.DownloadingItem.Pending {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}

	movie.GetLog().Info("Launching manual movie download")

	scheduler.Download(&movie)

	c.JSON(http.StatusOK, gin.H{})
	return
}

func abortMovieDownload(c *gin.Context) {
	id := c.Param("id")
	var movie Movie
	req := db.Client.Unscoped().Find(&movie, id)
	if err := req.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	downloader.AbortDownload(&movie)

	c.AbortWithStatus(http.StatusNoContent)
}

func updateMovie(c *gin.Context) {
	id := c.Param("id")
	var movie Movie
	req := db.Client.Find(&movie, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var movieFromRequest Movie
	c.Bind(&movieFromRequest)

	movie = movieFromRequest
	//db.Client.Save(&movie)

	c.JSON(http.StatusOK, movie)
}

func changeMovieDownloadedState(c *gin.Context) {
	id := c.Param("id")
	var movie Movie
	req := db.Client.Find(&movie, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var downloadingItemFromRequest DownloadingItem
	c.BindJSON(&downloadingItemFromRequest)

	movie.DownloadingItem.Downloaded = downloadingItemFromRequest.Downloaded
	db.Client.Save(&movie)

	if movie.DownloadingItem.Downloaded {
		stats.Stats.Movies.Downloaded += 1
	} else {
		stats.Stats.Movies.Downloaded -= 1
	}

	c.JSON(http.StatusOK, movie)
}

func changeMovieCustomTitle(c *gin.Context) {
	id := c.Param("id")
	var movie Movie
	req := db.Client.Find(&movie, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var movieFromRequest Movie
	c.BindJSON(&movieFromRequest)

	movie.CustomTitle = movieFromRequest.CustomTitle
	db.Client.Save(&movie)

	c.JSON(http.StatusOK, movie)
}

func useMovieDefaultTitle(c *gin.Context) {
	id := c.Param("id")
	var movie Movie
	req := db.Client.Find(&movie, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var movieFromRequest Movie
	c.BindJSON(&movieFromRequest)

	movie.UseDefaultTitle = movieFromRequest.UseDefaultTitle
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

	stats.Stats.Movies.Tracked += 1
	stats.Stats.Movies.Removed -= 1

	c.JSON(http.StatusOK, gin.H{})
}
