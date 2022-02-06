package mongodb

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	log     = logrus.New()
	Session *mongo.Client
)

func init() {
	log.Formatter = new(logrus.JSONFormatter)
}

func ConnectToMongo() {
	MONGO_URI := os.Getenv("MONGO_URI")
	if len(MONGO_URI) == 0 {
		MONGO_URI = "mongodb://localhost/auth-db"
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Panic("Fail connect to Mongo", err)
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Panic("Fail connect to Mongo", err)
		panic(err)
	}

	log.Info("Success connect to MongoDB")
	Session = client
}
