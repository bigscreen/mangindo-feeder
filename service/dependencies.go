package service

import "github.com/bigscreen/mangindo-feeder/client"

type Dependencies struct {
	MangaService   MangaService
	ChapterService ChapterService
}

func InstantiateDependencies() Dependencies {
	mangaClient := client.NewMangaClient()
	chapterClient := client.NewChapterClient()

	mangaService := NewMangaService(mangaClient)
	chapterService := NewChapterService(chapterClient)

	return Dependencies{
		MangaService:   mangaService,
		ChapterService: chapterService,
	}
}
