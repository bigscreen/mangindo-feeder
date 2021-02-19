package manager

import (
	"encoding/json"
	"errors"

	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/common"
	"github.com/bigscreen/mangindo-feeder/domain"
)

type contentCacheManager struct {
	cClient client.ContentClient
	cCache  cache.ContentCache
}

type ContentCacheManager interface {
	SetCache(titleID string, chapter float32) error
	GetCache(titleID string, chapter float32) (*domain.ContentListResponse, error)
}

func (m *contentCacheManager) SetCache(titleID string, chapter float32) error {
	cl, err := m.cClient.GetContentList(titleID, chapter)
	if err != nil {
		return err
	}

	cs, _ := json.Marshal(cl)

	return m.cCache.Set(titleID, common.GetFormattedChapterNumber(chapter), string(cs))
}

func (m *contentCacheManager) GetCache(titleID string, chapter float32) (*domain.ContentListResponse, error) {
	cs, err := m.cCache.Get(titleID, common.GetFormattedChapterNumber(chapter))
	if err != nil {
		return nil, err
	}

	var cl *domain.ContentListResponse
	err = json.Unmarshal([]byte(cs), &cl)
	if err != nil {
		return nil, errors.New("invalid content cache")
	}

	return cl, nil
}

func NewContentCacheManager(client client.ContentClient, cache cache.ContentCache) *contentCacheManager {
	return &contentCacheManager{
		cClient: client,
		cCache:  cache,
	}
}
