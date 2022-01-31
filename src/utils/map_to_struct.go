package utils

import (
	"encoding/json"
)

func MapToStruct(data interface{}, v interface{}) error {
	jsonData, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		LogError("error while marshaling json", marshalErr)
		return marshalErr
	}
	unmarshalErr := json.Unmarshal(jsonData, v)
	if unmarshalErr != nil {
		LogError("error while unmarshaling json", unmarshalErr)
		return unmarshalErr
	}
	return nil
}
