package notifications

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type SimpelNotificationData struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	ImageUrl string `json:"image_url"`
}

type FCMNotification struct {
	client *messaging.Client
}

func (f *FCMNotification) InitializeFirebase() {
	ctx := context.Background()
	app, err := getFirebaseApp()
	if err != nil {
		log.Fatal(err)
	}

	client, err2 := app.Messaging(ctx)
	if err2 != nil {
		log.Fatal(err2)
	}
	f.client = client
}

func (f *FCMNotification) SendSimpelNotification(fcmToken string, data SimpelNotificationData, notificationData map[string]string) {
	ctx := context.Background()
	message := &messaging.Message{
		Token: fcmToken,
		Notification: &messaging.Notification{
			Title:    data.Title,
			Body:     data.Body,
			ImageURL: data.ImageUrl,
		},
		Data: notificationData,
	}
	response, err := f.client.Send(ctx, message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully sent message: %s\n", response)
}

func getFirebaseApp() (*firebase.App, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("paperly-2f24f-firebase-adminsdk-qc22i-1f43991742.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, nil
}
