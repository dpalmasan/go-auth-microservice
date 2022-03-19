package session

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-auth-microservice/api/user"
	"github.com/go-auth-microservice/models"
	"github.com/go-auth-microservice/types"
	"github.com/go-auth-microservice/utils/crypto"
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

	router.Post("/jwt", func(w http.ResponseWriter, r *http.Request) {
		Login(userModel, w, r)
	})

	router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		Registration(userModel, w, r)
	})

	router.Post("/refresh", func(w http.ResponseWriter, r *http.Request) {
		Refresh(userModel, w, r)
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

	dbUser, err := userModel.GetByEmail(user.Email)
	if err != nil {
		log.Error("email not registered in our records")
		return
	}

	validPassword := crypto.CheckPasswordHash(user.Passwordhash, dbUser.Passwordhash)
	if !validPassword {
		log.Error("Incorrect password for the provided email")
		return
	}

	tokenString, refreshToken, err := CreateJWTToken(dbUser)
	if err != nil {
		log.Fatal(err)
		return
	}

	// This is based on JWT duration time. Maybe should be added in a config
	timeTTL := time.Minute * 5
	timeDuration := time.Now().Add(timeTTL)
	cookie := http.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		HttpOnly: true,
		Expires:  timeDuration,
	}

	http.SetCookie(w, &cookie)
	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
	}
	http.SetCookie(w, &refreshTokenCookie)

	userIdCookie := http.Cookie{
		Name:     "user_id",
		Value:    dbUser.Id,
		HttpOnly: true,
	}

	http.SetCookie(w, &userIdCookie)
	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{
		"tokens": {
			"access": "` + tokenString + `",
			"refresh": "` + refreshToken + `"
		}
	}`))
}

func Registration(userModel models.UserModel, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Error(err)
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	user.Create(userModel, w, r)
}

func Refresh(userModel models.UserModel, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	refreshTokenCookie, err := r.Cookie("refresh_token")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Error("Not authorized")
		return
	}

	refreshToken := refreshTokenCookie.Value

	if refreshToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		log.Error("Not authorized")
		return
	}

	status, err := models.CheckRefreshToken(refreshToken)
	if err != nil || status != true {
		w.WriteHeader(http.StatusUnauthorized)
		log.Error(err)
		return
	}

	userIdCookie, err := r.Cookie("user_id")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Error("Error while reading cookie!")
		return
	}

	userId := userIdCookie.Value

	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		log.Error("Invalid user!")
		return
	}

	dbUser, err := userModel.GetById(userId)
	if err != nil {
		log.Errorf("User %s not found (%s)", userId, err)
		return
	}

	tokenString, refreshToken, err := CreateJWTToken(dbUser)
	if err != nil {
		log.Fatal(err)
		return
	}

	timeTTL := time.Minute * 5
	timeDuration := time.Now().Add(timeTTL)
	cookie := http.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		HttpOnly: true,
		Expires:  timeDuration,
	}

	http.SetCookie(w, &cookie)
	newRefreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
	}
	http.SetCookie(w, &newRefreshTokenCookie)

	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{
		"tokens": {
			"access": "` + tokenString + `",
			"refresh": "` + refreshToken + `"
		}
	}`))
}