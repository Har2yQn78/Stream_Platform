package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	MongoDb := os.Getenv("MONGODB_URI")

	if MongoDb == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}

	fmt.Println("MONGODB_URI:", MongoDb)

	clientOptions := options.Client().ApplyURI(MongoDb)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil
	}

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(collecationName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	databaseName := os.Getenv("DATABASE_NAME")

	fmt.Println("DATABASE_NAME:", databaseName)

	collection := Client.Database(databaseName).Collection(collecationName)

	if collection == nil {
		return nil
	}

	return collection
}
