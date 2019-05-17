package cache

import (
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ContentCacheTestSuite struct {
	suite.Suite
	c *contentCache
	k string
}

func (s *ContentCacheTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()

	s.k = generateContentCacheKey(contentTitleId, contentChapter)
}

func (s *ContentCacheTestSuite) SetupTest() {
	s.c = NewContentCache()
}

func TestContentCacheTestSuite(t *testing.T) {
	suite.Run(t, new(ContentCacheTestSuite))
}

const (
	contentTitleId = "foo"
	contentChapter = "69"
)

func (s *ContentCacheTestSuite) TestSet_ReturnsNilError_WhenValueStored() {
	value := "lorem ipsum"
	err := s.c.Set(contentTitleId, contentChapter, value)
	assert.Nil(s.T(), err)

	result, _ := s.c.redisClient.Get(s.k).Result()
	assert.Equal(s.T(), value, result)

	s.c.redisClient.Del(s.k)
}

func (s *ContentCacheTestSuite) TestGet_ReturnsError_WhenKeyIsMissing() {
	val, err := s.c.Get(contentTitleId, contentChapter)

	assert.Equal(s.T(), "", val)
	assert.Equal(s.T(), "redis: nil", err.Error())
}

func (s *ContentCacheTestSuite) TestGet_ReturnsValue_WhenKeyExists() {
	s.c.redisClient.Set(s.k, "lorem ipsum", 5*time.Second)
	val, err := s.c.Get(contentTitleId, contentChapter)

	assert.Equal(s.T(), "lorem ipsum", val)
	assert.Nil(s.T(), err)

	s.c.redisClient.Del(s.k)
}

func (s *ContentCacheTestSuite) TestDelete_WhenKeyIsMissing() {
	err := s.c.Delete(contentTitleId, contentChapter)

	assert.Nil(s.T(), err)
}

func (s *ContentCacheTestSuite) TestDelete_WhenKeyExists() {
	s.c.redisClient.Set(s.k, "lorem ipsum", 5*time.Second)
	err := s.c.Delete(contentTitleId, contentChapter)
	val, _ := s.c.redisClient.Get(s.k).Result()

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), val)
}
