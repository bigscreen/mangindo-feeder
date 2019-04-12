package domain

type Manga struct {
	Id           string `json:"id"`
	Title        string `json:"judul"`
	TitleId      string `json:"hidden_komik"`
	IconURL      string `json:"icon_komik"`
	LastChapter  string `json:"hiddenNewChapter"`
	ModifiedDate string `json:"lastModified"`
	Genre        string `json:"genre"`
	Alias        string `json:"nama_lain"`
	Author       string `json:"pengarang"`
	Status       string `json:"status"`
	PublishYear  string `json:"published"`
	Summary      string `json:"summary"`
}

type MangaListResponse struct {
	Mangas []Manga `json:"komik"`
}
