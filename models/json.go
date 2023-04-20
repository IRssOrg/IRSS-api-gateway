package models

import (
	"encoding/json"
)

func JsonArray2Slice(array []byte) ([]string, error) {
	var dataArray []interface{}
	if err := json.Unmarshal(array, &dataArray); err != nil {
		return nil, err
	}
	//dataArray := dataMap["data"].([]interface{})
	var slice []string
	for _, v := range dataArray {
		slice = append(slice, v.(string))
	}
	return slice, nil
}
