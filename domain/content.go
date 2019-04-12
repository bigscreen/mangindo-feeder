package domain

type Content struct {
	ImageURL string `json:"url"`
	Page     int    `json:"page"`
}

type ContentListResponse struct {
	Contents []Content `json:"chapter"`
}
