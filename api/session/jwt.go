package session

import (
	"time"

	"github.com/go-auth-microservice/models"
)

func CreateJWTToken() (string, string, error) {
	timeTTL := time.Minute * 5
	timeDuration := time.Now().Add(timeTTL).Unix()

	tokenString, err := models.NewAccessToken(timeDuration)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := models.NewRefreshToken(timeTTL)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshToken, nil
}
