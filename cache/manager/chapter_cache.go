package manager

import (
	"encoding/json"
	"errors"
	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/domain"
)

type chapterCacheManager struct {
	cClient client.ChapterClient
	cCache  cache.ChapterCache
}

type ChapterCacheManager interface {
	SetCache(titleId string) error
	GetCache(titleId string) (*domain.ChapterListResponse, error)
}

func (m *chapterCacheManager) SetCache(titleId string) error {
	cl, err := m.cClient.GetChapterList(titleId)
	if err != nil {
		return err
	}

	cs, _ := json.Marshal(cl)

	return m.cCache.Set(titleId, string(cs))
}

func (m *chapterCacheManager) GetCache(titleId string) (*domain.ChapterListResponse, error) {
	cs, err := m.cCache.Get(titleId)
	if err != nil {
		return nil, err
	}

	var cl *domain.ChapterListResponse
	err = json.Unmarshal([]byte(cs), &cl)
	if err != nil {
		return nil, errors.New("invalid chapter cache")
	}

	return cl, nil
}

func NewChapterCacheManager(client client.ChapterClient, cache cache.ChapterCache) *chapterCacheManager {
	return &chapterCacheManager{
		cClient: client,
		cCache:  cache,
	}
}
