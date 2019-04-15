package handler

import (
	"encoding/json"
	"github.com/bigscreen/mangindo-feeder/contract"
	"net/http"
)

func respondWith(statusCode int, r *http.Request, w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if response != nil {
		json.NewEncoder(w).Encode(response)
		return
	}
}

func getErrorResponse(err error) contract.ErrorResponse {
	return contract.ErrorResponse{
		Success: false,
		Error:   err.Error(),
	}
}
