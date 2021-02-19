package domain

type Chapter struct {
	Number       float32 `json:"hidden_chapter"`
	Title        string  `json:"judul"`
	TitleID      string  `json:"hidden_komik"`
	ModifiedDate string  `json:"waktu"`
}

type ChapterListResponse struct {
	Chapters []Chapter `json:"komik"`
}
