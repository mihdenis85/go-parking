package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/amend-parking-backend/internal/config"
)

var Client *mongo.Client
var DB *mongo.Database

func InitializeDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.Settings.MongoDBURL))
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
		return err
	}

	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Error pinging MongoDB: %v", err)
		return err
	}

	DB = Client.Database(config.Settings.DBName)
	log.Println("Database initialized successfully.")
	return nil
}

func CloseDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return Client.Disconnect(ctx)
}
