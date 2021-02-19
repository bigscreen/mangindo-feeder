package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bigscreen/mangindo-feeder/contract"
)

func respondWith(statusCode int, r *http.Request, w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if response != nil {
		_ = json.NewEncoder(w).Encode(response)
		return
	}
}

func getErrorResponse(err error) contract.ErrorResponse {
	return contract.ErrorResponse{
		Success: false,
		Error:   err.Error(),
	}
}
