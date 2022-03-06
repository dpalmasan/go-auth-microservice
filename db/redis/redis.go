package redis

import (
	"github.com/go-auth-microservice/utils"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()

	Redis *redis.Client
)

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)
}

func ConnectToRedis() {
	// Get configuration
	REDIS_URL := utils.Getenv("REDIS_URL", "redis://localhost:6379/1")
	opt, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		panic(err)
	}
	Redis = redis.NewClient(opt)

	_, err = Redis.Ping().Result()
	if err != nil {
		log.Panicf("Cannot connect to REDIS_URL=%s", REDIS_URL)
		log.Panic(err)
		panic(err)
	}

	log.Info("Success connect to Redis")
}
