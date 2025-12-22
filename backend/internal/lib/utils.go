package lib

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(v any) error {
	json, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return fmt.Errorf("error marshalling to json: %w", err)
	}
	fmt.Println(string(json))
	return nil
}
