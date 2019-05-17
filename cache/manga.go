package cache

import (
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/go-redis/redis"
	"time"
)

type mangaCache struct {
	redisClient *redis.Client
}

type MangaCache interface {
	Set(value string) error
	Get() (string, error)
	Delete() error
}

const mangaCacheKey = "MangasCache"

func (c *mangaCache) Set(value string) error {
	err := c.redisClient.Set(mangaCacheKey, value, time.Duration(constants.MangaCacheExpirationInMn)*time.Minute).Err()
	if err != nil {
		logger.Errorf("Failed to set %s - %s", mangaCacheKey, err)
	}
	return err
}

func (c *mangaCache) Get() (string, error) {
	value, err := c.redisClient.Get(mangaCacheKey).Result()
	if err != nil {
		logger.Errorf("Failed to get %s - %s", mangaCacheKey, err)
	}
	return value, err
}

func (c *mangaCache) Delete() error {
	err := c.redisClient.Del(mangaCacheKey).Err()
	if err != nil {
		logger.Errorf("Failed to delete %s - %s", mangaCacheKey, err)
	}
	return err
}

func NewMangaCache() *mangaCache {
	return &mangaCache{
		redisClient: appcontext.GetRedisClient(),
	}
}
