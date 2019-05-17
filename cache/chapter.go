package cache

import (
	"fmt"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/go-redis/redis"
	"time"
)

type chapterCache struct {
	redisClient *redis.Client
}

type ChapterCache interface {
	Set(titleId, value string) error
	Get(titleId string) (string, error)
	Delete(titleId string) error
}

func generateChapterCacheKey(titleId string) string {
	return fmt.Sprintf("ChaptersCache|%s", titleId)
}

func (c *chapterCache) Set(titleId, value string) error {
	key := generateChapterCacheKey(titleId)
	err := c.redisClient.Set(key, value, time.Duration(constants.ChapterCacheExpirationInMn)*time.Minute).Err()
	if err != nil {
		logger.Errorf("Failed to set %s - %s", key, err)
	}
	return err
}

func (c *chapterCache) Get(titleId string) (string, error) {
	key := generateChapterCacheKey(titleId)
	value, err := c.redisClient.Get(key).Result()
	if err != nil {
		logger.Errorf("Failed to get %s - %s", key, err)
	}
	return value, err
}

func (c *chapterCache) Delete(titleId string) error {
	key := generateChapterCacheKey(titleId)
	err := c.redisClient.Del(key).Err()
	if err != nil {
		logger.Errorf("Failed to delete %s - %s", key, err)
	}
	return err
}

func NewChapterCache() *chapterCache {
	return &chapterCache{
		redisClient: appcontext.GetRedisClient(),
	}
}
