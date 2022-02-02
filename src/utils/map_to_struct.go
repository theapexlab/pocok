package utils

import (
	"encoding/json"
)

func MapToStruct(data interface{}, v interface{}) error {
	jsonData, marshalError := json.Marshal(data)
	if marshalError != nil {
		LogError("error while marshaling json", marshalError)
		return marshalError
	}
	unmarshalError := json.Unmarshal(jsonData, v)
	if unmarshalError != nil {
		LogError("error while unmarshaling json", unmarshalError)
		return unmarshalError
	}
	return nil
}
