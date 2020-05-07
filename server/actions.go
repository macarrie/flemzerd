package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macarrie/flemzerd/cache"
	"github.com/macarrie/flemzerd/db"
	"net/http"
	"os"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/scheduler"
)

func actionsCheckNow(c *gin.Context) {
	log.Info("Manual polling launch")
	scheduler.ResetRunTicker()
	c.JSON(http.StatusOK, gin.H{})
}

func resetAllMetadata(c *gin.Context) {
	log.Info("Refreshing all metadata")

	movieCacheFolderPath := fmt.Sprintf("%s/movies/", CACHE_BASE_PATH)
	if err := os.RemoveAll(movieCacheFolderPath); err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"path": movieCacheFolderPath,
		}).Error("Could not delete movie cache folder")
	} else {
		err := os.MkdirAll(movieCacheFolderPath, 0644)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"path": movieCacheFolderPath,
			}).Error("Could not recreate empty movies cache folder")
		}
	}

	showsCacheFolderPath := fmt.Sprintf("%s/shows/", CACHE_BASE_PATH)
	if err := os.RemoveAll(showsCacheFolderPath); err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"path": showsCacheFolderPath,
		}).Error("Could not delete shows cache folder")
	} else {
		err := os.MkdirAll(showsCacheFolderPath, 0644)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"path": showsCacheFolderPath,
			}).Error("Could not recreate empty shows cache folder")
		}
	}

	var movies []Movie
	db.Client.Find(&movies)
	for _, m := range movies {
		cache.Clear(&m)
		db.Client.Save(&m)
	}

	var shows []TvShow
	db.Client.Find(&shows)
	for _, s := range shows {
		cache.Clear(&s)
		db.Client.Save(&s)
	}

	scheduler.ResetRunTicker()

	c.AbortWithStatus(http.StatusCreated)
}
