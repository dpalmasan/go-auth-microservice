package mongodb

import (
	"context"

	"github.com/go-auth-microservice/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	log          = logrus.New()
	Session      *mongo.Client
	DatabaseName string
)

func init() {
	log.Formatter = new(logrus.JSONFormatter)
}

func ConnectToMongo() {
	DatabaseName = utils.Getenv("MONGO_DBNAME", "auth")
	MONGO_URI := utils.Getenv("MONGO_URI", "mongodb://localhost/auth")

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
