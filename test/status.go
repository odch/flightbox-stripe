package test

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
)

func UpdateStatus(config *Config, ctx context.Context, id string, status string) error {
	conf := &firebase.Config{
		DatabaseURL: config.FirebaseDatabaseUrl,
	}
	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	// As an admin, the app has access to read and write all data, regradless of Security Rules
	ref := client.NewRef("card-payments/" + id + "/status")
	//var data map[string]interface{}
	if err := ref.Set(ctx, &status); err != nil {
		log.Fatalln("Error writing  to database:", err)
	}
	return err
}

func UpdateData(config *Config, ctx context.Context, id string, data string) error {
	conf := &firebase.Config{
		DatabaseURL: config.FirebaseDatabaseUrl,
	}
	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	// As an admin, the app has access to read and write all data, regradless of Security Rules
	ref := client.NewRef("card-payments/" + id + "/data")
	//var data map[string]interface{}
	if err := ref.Set(ctx, &data); err != nil {
		log.Fatalln("Error writing  to database:", err)
	}
	return err
}
