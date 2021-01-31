package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	configVars := map[string]string{
		"APP_PORT":               "3001",
		"LOG_LEVEL":              "debug",
		"ENVIRONMENT":            "test",
		"REDIS_HOST":             "localhost",
		"REDIS_PORT":             "6379",
		"REDIS_POOL":             "10",
		"WORKER_REDIS_ADDRESS":   "127.0.0.1:6379",
		"ORIGIN_SERVER_BASE_URL": "https://foo.com",
		"POPULAR_MANGA_TAGS":     "foo1, foo2",
		"ADS_CONTENT_TAGS":       "foo1, foo2",
	}

	for k, v := range configVars {
		os.Setenv(k, v)
		defer os.Unsetenv(k)
	}

	Load()
	assert.Equal(t, 3001, Port())
	assert.Equal(t, configVars["LOG_LEVEL"], LogLevel())
	assert.Equal(t, configVars["REDIS_HOST"], RedisHost())
	assert.Equal(t, 6379, RedisPort())
	assert.Equal(t, 10, RedisPool())
	assert.Equal(t, configVars["WORKER_REDIS_ADDRESS"], WorkerRedisAddress())
	assert.Equal(t, configVars["ORIGIN_SERVER_BASE_URL"], BaseURL())
	assert.Equal(t, []string{"foo1", "foo2"}, PopularMangaTags())
	assert.Equal(t, []string{"foo1", "foo2"}, AdsContentTags())
}
