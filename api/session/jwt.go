package session

import (
	"time"

	"github.com/go-auth-microservice/models"
	"github.com/go-auth-microservice/types"
)

func CreateJWTToken(user types.User) (string, string, error) {
	timeTTL := time.Minute * 5
	timeDuration := time.Now().Add(timeTTL).Unix()

	tokenString, err := models.NewAccessToken(user, timeDuration)

	if err != nil {
		return "", "", err
	}

	timeTTL = time.Hour * 24 * 7
	refreshToken, err := models.NewRefreshToken(user, timeTTL)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshToken, nil
}
