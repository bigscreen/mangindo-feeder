package cache

import (
	"fmt"
	"time"

	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/go-redis/redis"
)

type contentCache struct {
	redisClient *redis.Client
}

type ContentCache interface {
	Set(titleID, chapter, value string) error
	Get(titleID, chapter string) (string, error)
	Delete(titleID, chapter string) error
}

func generateContentCacheKey(titleID, chapter string) string {
	return fmt.Sprintf("ContentsCache|%s|%s", titleID, chapter)
}

func (c *contentCache) Set(titleID, chapter, value string) error {
	key := generateContentCacheKey(titleID, chapter)
	err := c.redisClient.Set(key, value, time.Duration(constants.ContentCacheExpirationInMn)*time.Minute).Err()
	if err != nil {
		logger.Errorf("Failed to set %s - %s", key, err)
	}
	return err
}

func (c *contentCache) Get(titleID, chapter string) (string, error) {
	key := generateContentCacheKey(titleID, chapter)
	value, err := c.redisClient.Get(key).Result()
	if err != nil {
		logger.Errorf("Failed to get %s - %s", key, err)
	}
	return value, err
}

func (c *contentCache) Delete(titleID, chapter string) error {
	key := generateContentCacheKey(titleID, chapter)
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
