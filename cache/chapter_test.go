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

type ChapterCacheTestSuite struct {
	suite.Suite
	c *chapterCache
	k string
}

func (s *ChapterCacheTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()

	s.k = generateChapterCacheKey(chapterTitleId)
}

func (s *ChapterCacheTestSuite) SetupTest() {
	s.c = NewChapterCache()
}

func TestChapterCacheTestSuite(t *testing.T) {
	suite.Run(t, new(ChapterCacheTestSuite))
}

const (
	chapterTitleId = "foo"
)

func (s *ChapterCacheTestSuite) TestSet_ReturnsNilError_WhenValueStored() {
	value := "lorem ipsum"
	err := s.c.Set(chapterTitleId, value)
	assert.Nil(s.T(), err)

	result, _ := s.c.redisClient.Get(s.k).Result()
	assert.Equal(s.T(), value, result)

	s.c.redisClient.Del(s.k)
}

func (s *ChapterCacheTestSuite) TestGet_ReturnsError_WhenKeyIsMissing() {
	val, err := s.c.Get(chapterTitleId)

	assert.Equal(s.T(), "", val)
	assert.Equal(s.T(), "redis: nil", err.Error())
}

func (s *ChapterCacheTestSuite) TestGet_ReturnsValue_WhenKeyExists() {
	s.c.redisClient.Set(s.k, "lorem ipsum", 5*time.Second)
	val, err := s.c.Get(chapterTitleId)

	assert.Equal(s.T(), "lorem ipsum", val)
	assert.Nil(s.T(), err)

	s.c.redisClient.Del(s.k)
}

func (s *ChapterCacheTestSuite) TestDelete_WhenKeyIsMissing() {
	err := s.c.Delete(chapterTitleId)

	assert.Nil(s.T(), err)
}

func (s *ChapterCacheTestSuite) TestDelete_WhenKeyExists() {
	s.c.redisClient.Set(s.k, "lorem ipsum", 5*time.Second)
	err := s.c.Delete(chapterTitleId)
	val, _ := s.c.redisClient.Get(s.k).Result()

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), val)
}
