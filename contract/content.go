package contract

import "strconv"

type ContentRequest struct {
	TitleId string
	Chapter float32
}

type ContentResponse struct {
	BaseResponse
	Contents []Content `json:"contents,omitempty"`
}

type Content struct {
	ImageURL string `json:"image_url"`
}

func NewContentRequest(titleId, chapter string) ContentRequest {
	chapterF, err := strconv.ParseFloat(chapter, 32)
	if err != nil {
		chapterF = 0.0
	}

	return ContentRequest{
		TitleId: titleId,
		Chapter: float32(chapterF),
	}
}
