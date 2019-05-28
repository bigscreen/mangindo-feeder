package manager

import (
	"encoding/json"
	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/client"
)

type mangaCacheManager struct {
	mClient client.MangaClient
	mCache  cache.MangaCache
}

type MangaCacheManager interface {
	SetCache() error
}

func (m *mangaCacheManager) SetCache() error {
	ml, err := m.mClient.GetMangaList()
	if err != nil {
		return err
	}

	ms, _ := json.Marshal(ml)

	return m.mCache.Set(string(ms))
}

func NewMangaCacheManager(client client.MangaClient, cache cache.MangaCache) *mangaCacheManager {
	return &mangaCacheManager{
		mClient: client,
		mCache:  cache,
	}
}
