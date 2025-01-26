package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestApiResponse_SuccessResponse(t *testing.T) {
	apiResponse := &ApiResponse{}
	payload := map[string]string{"key": "value"}
	expectedResponse := map[string]interface{}{
		"status":  true,
		"message": "success",
		"payload": payload,
	}

	// Create a ResponseRecorder to capture the HTTP response
	rr := httptest.NewRecorder()

	// Call the SuccessResponse method
	apiResponse.SuccessResponse(rr, http.StatusOK, payload)

	// Verify the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Verify the response body
	var actualResponse map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(expectedResponse, actualResponse) {
		t.Errorf("Expected response %v, got %v", expectedResponse, actualResponse)
	}
}

func TestApiResponse_ErrorResponse(t *testing.T) {
	apiResponse := &ApiResponse{}
	payload := "Something went wrong"
	expectedResponse := map[string]interface{}{
		"status":  false,
		"message": "failed",
		"payload": payload,
	}

	// Create a ResponseRecorder to capture the HTTP response
	rr := httptest.NewRecorder()

	// Call the ErrorResponse method
	apiResponse.ErrorResponse(rr, payload)

	// Verify the status code (should be 400 for errors)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}

	// Verify the response body
	var actualResponse map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(expectedResponse, actualResponse) {
		t.Errorf("Expected response %v, got %v", expectedResponse, actualResponse)
	}
}





// func TestApiResponse_SuccessResponse(t *testing.T) {
// 	// Test cases
// 	tests := []struct {
// 		name       string
// 		statusCode int
// 		payload    interface{}
// 		want       map[string]interface{}
// 	}{
// 		{
// 			name:       "basic success response",
// 			statusCode: http.StatusOK,
// 			payload:    map[string]string{"data": "test"},
// 			want: map[string]interface{}{
// 				"status":  true,
// 				"message": "success",
// 				"payload": map[string]string{"data": "test"},
// 			},
// 		},
// 		{
// 			name:       "empty payload",
// 			statusCode: http.StatusCreated,
// 			payload:    nil,
// 			want: map[string]interface{}{
// 				"status":  true,
// 				"message": "success",
// 				"payload": nil,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Create a new response recorder (implements http.ResponseWriter)
// 			w := httptest.NewRecorder()
			
// 			// Create an instance of ApiResponse
// 			ar := &ApiResponse{}
			
// 			// Call the method we're testing
// 			ar.SuccessResponse(w, tt.statusCode, tt.payload)
			
// 			// Check status code
// 			if w.Code != tt.statusCode {
// 				t.Errorf("SuccessResponse() status code = %v, want %v", w.Code, tt.statusCode)
// 			}
			
// 			// Check Content-Type header
// 			if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
// 				t.Errorf("SuccessResponse() content-type = %v, want application/json", contentType)
// 			}
			
// 			// Decode the response body
// 			var got map[string]interface{}
// 			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
// 				t.Fatalf("Failed to decode response body: %v", err)
// 			}
			
// 			// Compare the response with expected values
// 			if got["status"] != tt.want["status"] {
// 				t.Errorf("SuccessResponse() status = %v, want %v", got["status"], tt.want["status"])
// 			}
// 			if got["message"] != tt.want["message"] {
// 				t.Errorf("SuccessResponse() message = %v, want %v", got["message"], tt.want["message"])
// 			}
// 		})
// 	}
// }

// func TestApiResponse_ErrorResponse(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		payload interface{}
// 		want    map[string]interface{}
// 	}{
// 		{
// 			name:    "basic error response",
// 			payload: "invalid input",
// 			want: map[string]interface{}{
// 				"status":  false,
// 				"message": "failed",
// 				"payload": "invalid input",
// 			},
// 		},
// 		{
// 			name:    "error with structured payload",
// 			payload: map[string]string{"error": "validation failed"},
// 			want: map[string]interface{}{
// 				"status":  false,
// 				"message": "failed",
// 				"payload": map[string]string{"error": "validation failed"},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			ar := &ApiResponse{}
			
// 			ar.ErrorResponse(w, tt.payload)
			
// 			// Error responses should always return 400
// 			if w.Code != http.StatusBadRequest {
// 				t.Errorf("ErrorResponse() status code = %v, want 400", w.Code)
// 			}
			
// 			if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
// 				t.Errorf("ErrorResponse() content-type = %v, want application/json", contentType)
// 			}
			
// 			var got map[string]interface{}
// 			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
// 				t.Fatalf("Failed to decode response body: %v", err)
// 			}
			
// 			if got["status"] != tt.want["status"] {
// 				t.Errorf("ErrorResponse() status = %v, want %v", got["status"], tt.want["status"])
// 			}
// 			if got["message"] != tt.want["message"] {
// 				t.Errorf("ErrorResponse() message = %v, want %v", got["message"], tt.want["message"])
// 			}
// 		})
// 	}
// }
