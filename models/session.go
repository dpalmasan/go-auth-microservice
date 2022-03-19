package models

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-auth-microservice/db/redis"
	"github.com/go-auth-microservice/types"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	PRIVATE_KEY = "cert/private_key.pem"
	PUBLIC_KEY  = "cert/public_key.pub"
)

var (
	log = logrus.New()

	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	log.Formatter = new(logrus.JSONFormatter)
	signBytes, err := ioutil.ReadFile(PRIVATE_KEY)

	if err != nil {
		log.Fatal(err)
		return
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
		return
	}

	verifyBytes, err := ioutil.ReadFile(PUBLIC_KEY)
	if err != nil {
		log.Fatal(err)
		return
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func NewAccessToken(user types.User, timeDuration int64) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = timeDuration
	claims["iat"] = time.Now().Unix()
	claims["role"] = user.Role.String()
	token.Claims = claims
	token.Header["kid"] = "4abaccf5-1b16-4ecf-aa98-c75091ffab5c"
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewRefreshToken(timeDuration time.Duration) (string, error) {
	refreshToken, _ := uuid.NewUUID()

	err := redis.Redis.Set(refreshToken.String(), "true", timeDuration).Err()
	if err != nil {
		return "", err
	}
	return refreshToken.String(), nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}

		return verifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, err
}

func CheckRefreshToken(token string) (bool, error) {
	value := redis.Redis.Get(token)
	if value.Err() != nil {
		return false, value.Err()
	}

	status, err := value.Result()
	if err != nil && status != "true" {
		return false, err
	}

	return true, nil
}