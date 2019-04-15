package contract

type ChapterRequest struct {
	TitleId string
}

type ChapterResponse struct {
	Success  bool      `json:"success"`
	Chapters []Chapter `json:"chapters"`
}

type Chapter struct {
	Number  string `json:"number"`
	Title   string `json:"title"`
	TitleId string `json:"title_id"`
}

func NewChapterRequest(titleId string) ChapterRequest {
	return ChapterRequest{TitleId: titleId}
}
