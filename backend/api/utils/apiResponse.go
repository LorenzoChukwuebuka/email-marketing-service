package utils

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
}

func (r *ApiResponse) SuccessResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	successResponse := map[string]interface{}{
		"status":    true,
		"message": "success",
		"payload": payload,
	}
	r.sendJSONResponse(w, statusCode, successResponse)
}

func (r *ApiResponse) ErrorResponse(w http.ResponseWriter, payload interface{}) {
	errorResponse := map[string]interface{}{
		"status":    false,
		"message": "failed",
		"payload": payload,
	}

	r.sendJSONResponse(w, 400, errorResponse)
}

func (r *ApiResponse) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
