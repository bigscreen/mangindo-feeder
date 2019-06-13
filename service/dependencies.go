package service

import (
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/cache/manager"
	"github.com/bigscreen/mangindo-feeder/client"
)

type Dependencies struct {
	MangaService   MangaService
	ChapterService ChapterService
	ContentService ContentService
}

type WorkerDependencies struct {
	MangaCacheManager   manager.MangaCacheManager
	ChapterCacheManager manager.ChapterCacheManager
	ContentCacheManager manager.ContentCacheManager
}

func InstantiateDependencies() Dependencies {
	mangaClient := client.NewMangaClient()
	chapterClient := client.NewChapterClient()
	contentClient := client.NewContentClient()

	mangaCache := cache.NewMangaCache()

	mangaCacheManager := manager.NewMangaCacheManager(mangaClient, mangaCache)

	workerService := NewWorkerService(appcontext.GetWorkerAdapter())

	mangaService := NewMangaService(mangaClient, mangaCacheManager, workerService)
	chapterService := NewChapterService(chapterClient)
	contentService := NewContentService(contentClient)

	return Dependencies{
		MangaService:   mangaService,
		ChapterService: chapterService,
		ContentService: contentService,
	}
}

func InstantiateWorkerDependencies() WorkerDependencies {
	mangaClient := client.NewMangaClient()
	chapterClient := client.NewChapterClient()
	contentClient := client.NewContentClient()

	mangaCache := cache.NewMangaCache()
	chapterCache := cache.NewChapterCache()
	contentCache := cache.NewContentCache()

	mangaCacheManager := manager.NewMangaCacheManager(mangaClient, mangaCache)
	chapterCacheManager := manager.NewChapterCacheManager(chapterClient, chapterCache)
	contentCacheManager := manager.NewContentCacheManager(contentClient, contentCache)

	return WorkerDependencies{
		MangaCacheManager:   mangaCacheManager,
		ChapterCacheManager: chapterCacheManager,
		ContentCacheManager: contentCacheManager,
	}
}
