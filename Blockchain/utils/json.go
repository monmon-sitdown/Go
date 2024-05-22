package utils

import (
	"encoding/json"
)

type StatusResponse struct {
	Status string `json:"status"`
}

func JsonStatus(status string) []byte {
	response := StatusResponse{Status: status}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// 如果序列化失败，返回一个通用的错误状态
		return []byte(`{"status":"error"}`)
	}
	return jsonResponse
}
