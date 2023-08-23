package utils

import (
	"encoding/json"
	"fmt"
)

func Output(v interface{}) {
	// Marshal the user struct into JSON
	jsonData, err := json.Marshal(v)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the JSON data
	fmt.Println(string(jsonData))

}
