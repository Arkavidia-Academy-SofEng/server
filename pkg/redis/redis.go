package redis

import (
	"context"
	"errors"
	redisPkg "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type ItfRedis interface {
	SetOTP(c context.Context, email string, code string) error
	GetOTP(c context.Context, email string) (string, error)
}

type redis struct {
	client *redisPkg.Client
}

func New() ItfRedis {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		logrus.Info("Failed to convert REDIS_DB to int")
	} else {
		logrus.Info("Successfully starting Redis")
	}

	client := redisPkg.NewClient(&redisPkg.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	return &redis{client: client}
}

func (r *redis) SetOTP(c context.Context, email string, code string) error {
	err := r.client.Set(c, email, code, 2*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redis) GetOTP(c context.Context, email string) (string, error) {
	val, err := r.client.Get(c, email).Result()
	if errors.Is(err, redisPkg.Nil) {
		return "", nil
	} else if err != nil {

		return "", err
	}
	return val, nil
}
