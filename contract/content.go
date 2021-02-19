package contract

import "strconv"

type ContentRequest struct {
	TitleID string
	Chapter float32
}

type ContentResponse struct {
	Success  bool      `json:"success"`
	Contents []Content `json:"contents"`
}

type Content struct {
	ImageURL string `json:"image_url"`
}

func NewContentRequest(titleID, chapter string) ContentRequest {
	chapterF, err := strconv.ParseFloat(chapter, 32)
	if err != nil {
		chapterF = 0.0
	}

	return ContentRequest{
		TitleID: titleID,
		Chapter: float32(chapterF),
	}
}
