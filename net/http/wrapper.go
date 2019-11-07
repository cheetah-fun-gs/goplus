package http

import (
	"time"
)

func defaultClient() *Client {
	return &Client{
		Timeout: time.Second * 5,
	}
}

// JSON post 方式获取 json返回
func JSON(url string, request interface{}, response interface{}) error {
	return defaultClient().JSON(url, request, response)
}

// GetJSON 用get方式获取json返回
func GetJSON(url string, response interface{}) error {
	return defaultClient().GetJSON(url, response)
}

// PostJSON 用post方式获取json返回
func PostJSON(url string, response interface{}) error {
	return defaultClient().PostJSON(url, response)
}

// JSONMultipartForm 上传文件或其他表单数据
func JSONMultipartForm(url string, fields []*MultipartFormField, response interface{}) error {
	return defaultClient().JSONMultipartForm(url, fields, response)
}
