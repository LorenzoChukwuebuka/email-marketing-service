package helper

import (
	"encoding/json"
	"net/http"
)

func SuccessResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	successResponse := map[string]interface{}{
		"status":  true,
		"message": "success",
		"payload": payload,
	}
	sendJSONResponse(w, statusCode, successResponse)
}

// ErrorResponse with error message
func ErrorResponse(w http.ResponseWriter, err error, payload interface{}) {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	} else {
		errorMessage = "An error occurred"
	}
	errorResponse := map[string]interface{}{
		"status":  false,
		"message": "failed",
		"error":   errorMessage, // Include the actual error message
		"payload": payload,
	}
	sendJSONResponse(w, http.StatusBadRequest, errorResponse)
}

// ErrorResponseWithCode allows setting custom status code
func ErrorResponseWithCode(w http.ResponseWriter, statusCode int, err error, payload interface{}) {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	} else {
		errorMessage = "An error occurred"
	}
	
	errorResponse := map[string]interface{}{
		"status":  false,
		"message": "failed",
		"error":   errorMessage, // Include the actual error message
		"payload": payload,
	}

	sendJSONResponse(w, statusCode, errorResponse)
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}