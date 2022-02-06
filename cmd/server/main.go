package main

import (
	"context"
	"time"

	"github.com/go-auth-microservice/db/mongodb"
	"github.com/go-auth-microservice/models"
	"github.com/go-auth-microservice/models/providers"
	"github.com/go-auth-microservice/types"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)

	// Connect to DB
	mongodb.ConnectToMongo()
}

func main() {
	var db models.UserModel

	db = providers.MongoDBUser{}

	defer func() {
		if err := mongodb.Session.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	user := types.User{
		Email:        "test-user@gotest.cl",
		Username:     "user1",
		Passwordhash: "11123451",
		CreatedAt:    time.Date(2020, 11, 14, 11, 30, 32, 0, time.UTC),
		Role:         0,
	}

	_, err := db.Add(user)
	if err != nil {
		log.Error(err)
	} else {
		log.Info("User inserted successfully!")
	}

	user, err = db.GetByEmail("test-user@gotest.cl")

	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("Retrieved user %+v", user)
}
