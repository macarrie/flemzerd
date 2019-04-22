package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/stats"

	. "github.com/macarrie/flemzerd/objects"
)

func getNotifications(c *gin.Context) {
	notifs, err := db.GetNotifications()
	if err != nil {
		log.Error("Error while getting notifications from db: ", err)
	}

	c.JSON(http.StatusOK, notifs)
}

func deleteNotifications(c *gin.Context) {
	db.Client.Unscoped().Delete(&Notification{})

	stats.Stats.Notifications.Read = 0
	stats.Stats.Notifications.Unread = 0

	c.JSON(http.StatusOK, gin.H{})
}

func getReadNotifications(c *gin.Context) {
	notifs, err := db.GetReadNotifications()
	if err != nil {
		log.Error("Error while getting read notifications from db: ", err)
	}

	c.JSON(http.StatusOK, notifs)
}

func markAllNotificationsAsRead(c *gin.Context) {
	db.Client.Model(Notification{}).Updates(Notification{Read: true})

	stats.Stats.Notifications.Read += stats.Stats.Notifications.Unread
	stats.Stats.Notifications.Unread = 0

	c.JSON(http.StatusOK, gin.H{})
}

func getUnreadNotifications(c *gin.Context) {
	notifs, err := db.GetUnreadNotifications()
	if err != nil {
		log.Error("Error while getting notifications from db: ", err)
	}

	c.JSON(http.StatusOK, notifs)
}

func changeNotificationReadState(c *gin.Context) {
	id := c.Param("id")
	var notif Notification
	req := db.Client.Find(&notif, id)
	if req.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var notifInfo Notification
	c.BindJSON(&notifInfo)

	notif.Read = notifInfo.Read
	db.Client.Save(&notif)

	if notif.Read {
		stats.Stats.Notifications.Read += 1
		stats.Stats.Notifications.Unread -= 1
	} else {
		stats.Stats.Notifications.Read -= 1
		stats.Stats.Notifications.Unread += 1
	}

	c.JSON(http.StatusOK, notif)
}

func renderNotificationsList(c *gin.Context) {
	var notifs []Notification
	db.Client.Order("created_at DESC").Find(&notifs)

	//Do not use master template for partial renders
	c.HTML(http.StatusOK, "notifications_list_render.tpl", gin.H{
		"notifications": notifs,
	})
}
