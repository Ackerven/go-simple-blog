package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Map map[string]interface{}

// 将request中json格式字符串转换为map类型

func RequestJsonInterface(r *http.Request) map[string]interface{} {

	var data map[string]interface{}

	//读取request body
	rawData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("读取失败")
	}

	//将rewData反序列化进map中
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		fmt.Println("反序列化失败")
	}
	return data
}