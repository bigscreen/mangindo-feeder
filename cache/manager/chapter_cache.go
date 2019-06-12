package manager

import (
	"encoding/json"
	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/client"
)

type chapterCacheManager struct {
	cClient client.ChapterClient
	cCache  cache.ChapterCache
}

type ChapterCacheManager interface {
	SetCache(titleId string) error
}

func (m *chapterCacheManager) SetCache(titleId string) error {
	cl, err := m.cClient.GetChapterList(titleId)
	if err != nil {
		return err
	}

	cs, _ := json.Marshal(cl)

	return m.cCache.Set(titleId, string(cs))
}

func NewChapterCacheManager(client client.ChapterClient, cache cache.ChapterCache) *chapterCacheManager {
	return &chapterCacheManager{
		cClient: client,
		cCache:  cache,
	}
}
