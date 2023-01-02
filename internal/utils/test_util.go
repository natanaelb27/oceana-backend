package utils

import (
	"bytes"
	"encoding/json"
)

func NewBody(data any) *bytes.Buffer {
	json, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return bytes.NewBuffer(json)
}
