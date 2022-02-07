package session

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-auth-microservice/models"
	"github.com/go-auth-microservice/types"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Formatter = new(logrus.JSONFormatter)
}

func Routes(userModel models.UserModel) chi.Router {
	router := chi.NewRouter()

	router.Use(chiMiddleware.AllowContentType("application/json"))

	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		Login(userModel, w, r)
	})

	return router
}

func Login(userModel models.UserModel, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Error(err)
		return
	}

	var user *types.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Error(err)
		return
	}

	password := user.Passwordhash
	dbUser, err := userModel.GetByEmail(user.Email)
	if err != nil {
		log.Error("email not registered in our records")
		return
	}

	validPassword := dbUser.ValidatePasswordHash(password)
	if !validPassword {
		log.Error("Incorrect password for the provided email")
		return
	}

	tokenString, refreshToken, err := CreateJWTToken()
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{
		"tokens": {
			"access": "` + tokenString + `",
			"refresh": "` + refreshToken + `"
		}
	}`))
}
