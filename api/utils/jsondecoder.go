package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeRequestBody(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func EncodeToJson(v interface{}) {
	// Marshal the user struct into JSON
	jsonData, err := json.Marshal(v)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the JSON data
	fmt.Println(string(jsonData))

}
