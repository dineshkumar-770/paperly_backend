package controller

import (
	"mongo_api/notifications"
	"net/http"
)

type NotificationController struct {
}

func (nc *NotificationController) SendSimpleNotification(w http.ResponseWriter, r *http.Request) {

	fcmToken := r.FormValue("token")
	notiTitle := r.FormValue("title")
	notiBody := r.FormValue("body")
	notiImageUrl := r.FormValue("image_url")
	notificationHeaders := notifications.SimpelNotificationData{
		Title:    notiTitle,
		Body:     notiBody,
		ImageUrl: notiImageUrl,
	}

	notificationData := map[string]string{
		"title":     "New Wallpapers Added!",
		"body":      "Check out the latest collection of wallpapers in your favorite categories.",
		"category":  "Nature",
		"timestamp": "2024-10-11T10:15:30Z",
		"type":      "new_wallpapers",
		"image_url": notiImageUrl,
	}

	fcmInstanace := notifications.FCMNotification{}
	fcmInstanace.InitializeFirebase()
	fcmInstanace.SendSimpelNotification(fcmToken, notificationHeaders, notificationData)
}
