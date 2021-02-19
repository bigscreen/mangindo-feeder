package cache

import (
	"testing"
	"time"

	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MangaCacheTestSuite struct {
	suite.Suite
	c *mangaCache
}

func (s *MangaCacheTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func (s *MangaCacheTestSuite) SetupTest() {
	s.c = NewMangaCache()
}

func TestMangaCacheTestSuite(t *testing.T) {
	suite.Run(t, new(MangaCacheTestSuite))
}

func (s *MangaCacheTestSuite) TestSet_ReturnsNilError_WhenValueStored() {
	value := "lorem ipsum"
	err := s.c.Set(value)
	assert.Nil(s.T(), err)

	result, _ := s.c.redisClient.Get(mangaCacheKey).Result()
	assert.Equal(s.T(), value, result)

	s.c.redisClient.Del(mangaCacheKey)
}

func (s *MangaCacheTestSuite) TestGet_ReturnsError_WhenKeyIsMissing() {
	val, err := s.c.Get()

	assert.Equal(s.T(), "", val)
	assert.Equal(s.T(), "redis: nil", err.Error())
}

func (s *MangaCacheTestSuite) TestGet_ReturnsValue_WhenKeyExists() {
	s.c.redisClient.Set(mangaCacheKey, "lorem ipsum", 5*time.Second)
	val, err := s.c.Get()

	assert.Equal(s.T(), "lorem ipsum", val)
	assert.Nil(s.T(), err)

	s.c.redisClient.Del(mangaCacheKey)
}

func (s *MangaCacheTestSuite) TestDelete_WhenKeyIsMissing() {
	err := s.c.Delete()

	assert.Nil(s.T(), err)
}

func (s *MangaCacheTestSuite) TestDelete_WhenKeyExists() {
	s.c.redisClient.Set(mangaCacheKey, "lorem ipsum", 5*time.Second)
	err := s.c.Delete()
	val, _ := s.c.redisClient.Get(mangaCacheKey).Result()

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), val)
}
