package service

import "github.com/bigscreen/mangindo-feeder/client"

type Dependencies struct {
	MangaService MangaService
}

func InstantiateDependencies() Dependencies {
	mangaClient := client.NewMangaClient()

	mangaService := NewMangaService(mangaClient)

	return Dependencies{
		MangaService: mangaService,
	}
}
