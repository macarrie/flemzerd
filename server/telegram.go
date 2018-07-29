package server

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/flemzerd/logging"
	notifier "github.com/macarrie/flemzerd/notifiers"
	"github.com/macarrie/flemzerd/notifiers/impl/telegram"
)

func performTelegramAuth(c *gin.Context) {
	n, err := notifier.GetNotifier("telegram")
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	t := n.(*telegram.TelegramNotifier)
	if t.ChatID != 0 {
		log.Error("CHAT ID NOT ZERO. EXIT STATUS NO CONTENT")
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	go t.Auth()
	c.JSON(http.StatusOK, gin.H{})
	return
}

func getTelegramChatID(c *gin.Context) {
	n, err := notifier.GetNotifier("telegram")
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	t := n.(*telegram.TelegramNotifier)
	c.JSON(http.StatusOK, gin.H{
		"chat_id": t.ChatID,
	})
}

func getTelegramAuthCode(c *gin.Context) {
	n, err := notifier.GetNotifier("telegram")
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	t := n.(*telegram.TelegramNotifier)
	if t.AuthCode == 0 {
		t.AuthCode = rand.Intn(1000000)
	}

	c.JSON(http.StatusOK, gin.H{
		"auth_code": t.AuthCode,
	})
}