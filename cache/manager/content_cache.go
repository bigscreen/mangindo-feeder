package manager

import (
	"encoding/json"
	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/common"
)

type contentCacheManager struct {
	cClient client.ContentClient
	cCache  cache.ContentCache
}

type ContentCacheManager interface {
	SetCache(titleId string, chapter float32) error
}

func (m *contentCacheManager) SetCache(titleId string, chapter float32) error {
	cl, err := m.cClient.GetContentList(titleId, chapter)
	if err != nil {
		return err
	}

	cs, _ := json.Marshal(cl)

	return m.cCache.Set(titleId, common.GetFormattedChapterNumber(chapter), string(cs))
}

func NewContentCacheManager(client client.ContentClient, cache cache.ContentCache) *contentCacheManager {
	return &contentCacheManager{
		cClient: client,
		cCache:  cache,
	}
}
