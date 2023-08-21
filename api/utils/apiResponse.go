package utils

import (
	"encoding/json"
	"net/http"
)

func SuccessResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	successResponse := map[string]interface{}{
		"code":    "1",
		"message": "success",
		"payload": payload,
	}
	sendJSONResponse(w, statusCode, successResponse)
}

func ErrorResponse(w http.ResponseWriter, payload interface{}) {
	errorResponse := map[string]interface{}{
		"code":    "3",
		"message": "failed",
		"payload": payload,
	}

	sendJSONResponse(w, 400, errorResponse)
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
