package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/macarrie/flemzerd/db"

	. "github.com/macarrie/flemzerd/objects"
)

func getNotifications(c *gin.Context) {
	var notifs []Notification
	db.Client.Order("created_at DESC").Find(&notifs)

	c.JSON(http.StatusOK, notifs)
}

func deleteNotifications(c *gin.Context) {
	db.Client.Unscoped().Delete(&Notification{})

	c.JSON(http.StatusOK, gin.H{})
}

func getReadNotifications(c *gin.Context) {
	var notifs []Notification
	db.Client.Where(&Notification{Read: true}).Order("created_at DESC").Find(&notifs)

	c.JSON(http.StatusOK, notifs)
}

func markAllNotificationsAsRead(c *gin.Context) {
	db.Client.Model(Notification{}).Updates(Notification{Read: true})

	c.JSON(http.StatusOK, gin.H{})
}

func getUnreadNotifications(c *gin.Context) {
	var notifs []Notification
	// Use plain SQL string instead of struct because GORM does not perform where on zero values when querying with structs
	db.Client.Where("read = false").Find(&notifs)

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
