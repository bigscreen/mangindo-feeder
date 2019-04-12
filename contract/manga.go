package contract

type Manga struct {
	Title       string `json:"title"`
	TitleId     string `json:"title_id"`
	IconURL     string `json:"icon_url"`
	LastChapter string `json:"last_chapter"`
	Genre       string `json:"genre"`
	Alias       string `json:"alias"`
	Author      string `json:"author"`
	Status      string `json:"status"`
	PublishYear string `json:"publish_date"`
	Summary     string `json:"summary"`
}

type MangaResponse struct {
	BaseResponse
	PopularMangas []Manga `json:"popular_mangas,omitempty"`
	LatestMangas  []Manga `json:"latest_mangas,omitempty"`
}
