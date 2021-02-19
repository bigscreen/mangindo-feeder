package contract

type ChapterRequest struct {
	TitleID string
}

type Chapter struct {
	Number  string `json:"number"`
	Title   string `json:"title"`
	TitleID string `json:"title_id"`
}

type ChapterResponse struct {
	Success  bool      `json:"success"`
	Chapters []Chapter `json:"chapters"`
}

func NewChapterRequest(titleID string) ChapterRequest {
	return ChapterRequest{TitleID: titleID}
}
