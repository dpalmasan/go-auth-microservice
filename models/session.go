package models

import (
	"crypto/rsa"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func NewAccessToken(timeDuration int64) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = timeDuration
	claims["iat"] = time.Now().Unix()
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewRefreshToken(timeDuration time.Duration) (string, error) {
	refreshToken, _ := uuid.NewUUID()

	return refreshToken.String(), nil
}
