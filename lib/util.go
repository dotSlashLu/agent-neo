package lib

import (
	"bytes"
	"encoding/json"
)

func TrimBuf(buf []byte) []byte {
	return bytes.Trim(buf, "\x00")
}

func RespError(e error) ([]byte, error) {
	type resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	str, err := json.Marshal(resp{"error", e.Error()})
	if err != nil {
		panic(err)
	}
	return str, nil
}
