package domain

type Chapter struct {
	Number       int    `json:"hidden_chapter"`
	Title        string `json:"judul"`
	TitleId      string `json:"hidden_komik"`
	ModifiedDate string `json:"waktu"`
}

type ChapterListResponse struct {
	Chapters []Chapter `json:"komik"`
}