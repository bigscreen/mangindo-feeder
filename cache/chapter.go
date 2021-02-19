package cache

import (
	"fmt"
	"time"

	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/go-redis/redis"
)

type chapterCache struct {
	redisClient *redis.Client
}

type ChapterCache interface {
	Set(titleID, value string) error
	Get(titleID string) (string, error)
	Delete(titleID string) error
}

func generateChapterCacheKey(titleID string) string {
	return fmt.Sprintf("ChaptersCache|%s", titleID)
}

func (c *chapterCache) Set(titleID, value string) error {
	key := generateChapterCacheKey(titleID)
	err := c.redisClient.Set(key, value, time.Duration(constants.ChapterCacheExpirationInMn)*time.Minute).Err()
	if err != nil {
		logger.Errorf("Failed to set %s - %s", key, err)
	}
	return err
}

func (c *chapterCache) Get(titleID string) (string, error) {
	key := generateChapterCacheKey(titleID)
	value, err := c.redisClient.Get(key).Result()
	if err != nil {
		logger.Errorf("Failed to get %s - %s", key, err)
	}
	return value, err
}

func (c *chapterCache) Delete(titleID string) error {
	key := generateChapterCacheKey(titleID)
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
