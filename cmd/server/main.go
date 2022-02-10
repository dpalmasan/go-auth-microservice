package main

import (
	"context"
	"net/http"

	"github.com/go-auth-microservice/api/session"
	"github.com/go-auth-microservice/api/user"
	"github.com/go-auth-microservice/db/mongodb"
	"github.com/go-auth-microservice/models/providers"
	"github.com/go-chi/chi"
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
	PORT := "4000"
	db := providers.MongoDBUser{}

	defer func() {
		if err := mongodb.Session.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	router := chi.NewRouter()

	router.Mount("/login", session.Routes(db))
	router.Mount("/users", user.Routes(db))

	log.Infof("Running service on port %s", PORT)
	http.ListenAndServe(":"+PORT, router)
}
