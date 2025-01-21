package utils

import (
	"encoding/json"
	"net/http"
)

func DecodeRequestBody(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func EncodeToJson(v interface{}) (map[string]interface{}, error) {
	// Marshal the input value into JSON
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Create a map to hold the JSON data
	var result map[string]interface{}

	// Unmarshal the JSON data into the result map
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	return result, nil
}
