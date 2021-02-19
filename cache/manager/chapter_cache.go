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
	SetCache(titleID string) error
	GetCache(titleID string) (*domain.ChapterListResponse, error)
}

func (m *chapterCacheManager) SetCache(titleID string) error {
	cl, err := m.cClient.GetChapterList(titleID)
	if err != nil {
		return err
	}

	cs, _ := json.Marshal(cl)

	return m.cCache.Set(titleID, string(cs))
}

func (m *chapterCacheManager) GetCache(titleID string) (*domain.ChapterListResponse, error) {
	cs, err := m.cCache.Get(titleID)
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
