package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/scheduler"
)

func actionsCheckNow(c *gin.Context) {
	log.Info("Manual polling launch")
	scheduler.ResetRunTicker()
	c.JSON(http.StatusOK, gin.H{})
}
