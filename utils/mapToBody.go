package utils

import (
	"encoding/json"
	"fmt"
)

// 将map序列化为[]byte

func MapToBody(data map[string]interface{}) []byte {
	v, err := json.Marshal(data)
	if err != nil {
		fmt.Println("序列化失败")
	}
	return v
}