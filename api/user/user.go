package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-auth-microservice/models"
	"github.com/go-auth-microservice/types"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Formatter = new(logrus.JSONFormatter)
}

func Routes(userModel models.UserModel) chi.Router {
	router := chi.NewRouter()

	router.Use(chiMiddleware.AllowContentType("application/json"))

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		Create(userModel, w, r)
	})

	return router
}

func Create(userModel models.UserModel, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parent := opentracing.GlobalTracer().StartSpan("POST /users")

	defer parent.Finish()

	b, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		log.Error(err)
		return
	}

	var user types.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Error(err)
		return
	}

	// TODO: Implement a role to be passed so we overwrite this
	// By default all users should be regular!
	user.Role = types.Regular
	user, err = userModel.Add(user)

	if err != nil {
		log.Error(err)
		return
	}

	output, err := json.Marshal(user)
	if err != nil {
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(output)
}
