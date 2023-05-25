package lib

import (
	"encoding/json"
)

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
