package lib

import (
	"encoding/json"
	"fmt"
)

func AppendIfNotEmpty(payload *string, key, value string) {
	if len(key) == 0 || len(value) == 0 {
		return
	}

	if len(*payload) > 0 {
		*payload += ","
	}

	*payload += fmt.Sprintf(`"%s":"%s"`, key, value)
}

func DecodeFromJSON(in []byte, out any) error {
	return json.Unmarshal(in, out)
}

func EncodeToJSON(in any) ([]byte, error) {
	return json.Marshal(in)
}

func FormatTopLevelJsonArray(rawJson []any) []byte {
	data, err := json.Marshal(rawJson)

	if err != nil {
		return nil
	}

	return data
}

func FormatTopLevelJsonObject(rawJson map[string]any) []byte {
	data, err := json.Marshal(rawJson)

	if err != nil {
		return nil
	}

	return data
}

func ParseJsonArray(value any) []any {
	arrValue, ok := value.([]any)

	if ok {
		return arrValue
	}

	return nil
}

func ParseJsonBool(value any) bool {
	boolValue, ok := value.(bool)

	if ok {
		return boolValue
	}

	return false
}

func ParseJsonNumber(value any) float64 {
	numValue, ok := value.(float64)

	if ok {
		return numValue
	}

	return 0
}

func ParseJsonObject(value any) map[string]any {
	objValue, ok := value.(map[string]any)

	if ok {
		return objValue
	}

	return nil
}

func ParseJsonString(value any) string {
	strValue, ok := value.(string)

	if ok {
		return strValue
	}

	return ""
}

func ParseJsonStringArray(value any) []string {
	arrValue, ok := value.([]any)

	if ok {
		strArrValue := make([]string, len(arrValue))

		for idx, value := range arrValue {
			strArrValue[idx], ok = value.(string)

			if !ok {
				return nil
			}
		}

		return strArrValue
	}

	return nil
}

func ParseTopLevelJsonArray(data []byte) []any {
	var arrValue []any

	err := json.Unmarshal(data, &arrValue)

	if err != nil {
		return nil
	}

	return arrValue
}

func ParseTopLevelJsonObject(data []byte) map[string]any {
	var objValue map[string]any

	err := json.Unmarshal(data, &objValue)

	if err != nil {
		return nil
	}

	return objValue
}
