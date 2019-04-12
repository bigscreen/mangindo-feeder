package contract

type BaseResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
