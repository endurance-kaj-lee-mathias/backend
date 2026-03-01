package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func initializeFirebase() *firebase.App {
	opt := option.WithAuthCredentialsFile(option.ServiceAccount, "../../endurance-credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing app: %v\n", err)
	}
	return app
}
