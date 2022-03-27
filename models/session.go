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
	"github.com/sirupsen/logrus"
)

const (
	ACCESS_TOKEN_PRIVATE_KEY = "cert/private_key.pem"
	ACCESS_TOKEN_PUBLIC_KEY  = "cert/public_key.pub"
	REFRESH_TOKEN_PRIVATE_KEY = "cert/refresh_private_key.pem"
	REFRESH_TOKEN_PUBLIC_KEY  = "cert/refresh_public_key.pub"
)

var (
	log = logrus.New()

	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
	refreshVerifyKey *rsa.PublicKey
	refreshSignKey   *rsa.PrivateKey
)

func init() {
	log.Formatter = new(logrus.JSONFormatter)
	signBytes, err := ioutil.ReadFile(ACCESS_TOKEN_PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
		return
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
		return
	}

	verifyBytes, err := ioutil.ReadFile(ACCESS_TOKEN_PUBLIC_KEY)
	if err != nil {
		log.Fatal(err)
		return
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
		return
	}

	refreshSignBytes, err := ioutil.ReadFile(REFRESH_TOKEN_PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
		return
	}

	refreshSignKey, err = jwt.ParseRSAPrivateKeyFromPEM(refreshSignBytes)
	if err != nil {
		log.Fatal(err)
		return
	}

	refreshVerifyBytes, err := ioutil.ReadFile(REFRESH_TOKEN_PUBLIC_KEY)
	if err != nil {
		log.Fatal(err)
		return
	}

	refreshVerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(refreshVerifyBytes)
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
	claims["user_id"] = user.Id
	token.Claims = claims
	token.Header["kid"] = "4abaccf5-1b16-4ecf-aa98-c75091ffab5c"
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewRefreshToken(user types.User, timeDuration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = timeDuration
	claims["iat"] = time.Now().Unix()
	claims["role"] = user.Role.String()
	claims["user_id"] = user.Id
	token.Claims = claims
	token.Header["kid"] = "4abaccf5-1b16-4ecf-aa98-c75091ffab5c"
	tokenString, err := token.SignedString(refreshSignKey)
	if err != nil {
		return "", err
	}

	err = redis.Redis.Set(tokenString, "true", timeDuration).Err()
	if err != nil {
		return "", err
	}
	return tokenString, nil
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

func VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}
		return refreshVerifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, err
}

func CheckRefreshToken(token string) (*jwt.Token, error) {
	value := redis.Redis.Get(token)
	if value.Err() != nil {
		return nil, value.Err()
	}

	status, err := value.Result()
	if err != nil && status != "true" {
		return nil, err
	}

	refreshToken, err := VerifyRefreshToken(token)
	if err != nil {
		return nil, err
	}
	if _, ok := refreshToken.Claims.(jwt.Claims); !ok && !refreshToken.Valid {
		return nil, err
	}

	return refreshToken, nil
}

func RevokeRefreshToken(token string) {
	redis.Redis.Del(token)
}