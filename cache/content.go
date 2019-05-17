package cache

import (
	"fmt"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/go-redis/redis"
	"time"
)

type contentCache struct {
	redisClient *redis.Client
}

type ContentCache interface {
	Set(titleId, chapter, value string) error
	Get(titleId, chapter string) (string, error)
	Delete(titleId, chapter string) error
}

func generateContentCacheKey(titleId, chapter string) string {
	return fmt.Sprintf("ContentsCache|%s|%s", titleId, chapter)
}

func (c *contentCache) Set(titleId, chapter, value string) error {
	key := generateContentCacheKey(titleId, chapter)
	err := c.redisClient.Set(key, value, time.Duration(constants.ContentCacheExpirationInMn)*time.Minute).Err()
	if err != nil {
		logger.Errorf("Failed to set %s - %s", key, err)
	}
	return err
}

func (c *contentCache) Get(titleId, chapter string) (string, error) {
	key := generateContentCacheKey(titleId, chapter)
	value, err := c.redisClient.Get(key).Result()
	if err != nil {
		logger.Errorf("Failed to get %s - %s", key, err)
	}
	return value, err
}

func (c *contentCache) Delete(titleId, chapter string) error {
	key := generateContentCacheKey(titleId, chapter)
	err := c.redisClient.Del(key).Err()
	if err != nil {
		logger.Errorf("Failed to delete %s - %s", key, err)
	}
	return err
}

func NewContentCache() *contentCache {
	return &contentCache{
		redisClient: appcontext.GetRedisClient(),
	}
}
