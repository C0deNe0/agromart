package lib

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(v interface{}) {
	json, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Errorf("error marshalling to json: %w", err)
		return
	}
	fmt.Println(string(json)) 
}
