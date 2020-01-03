package http

import (
	"time"
)

func defaultClient() *Client {
	return &Client{
		Timeout: time.Second * 5,
	}
}

// JSON method:post req:json resp:json
func JSON(toURL string, request interface{}, response interface{}) error {
	return defaultClient().JSON(toURL, request, response)
}

// GetJSON method:get resp:json
func GetJSON(toURL string, response interface{}) error {
	return defaultClient().GetJSON(toURL, response)
}

// PostJSON method:get resp:json
func PostJSON(toURL string, response interface{}) error {
	return defaultClient().PostJSON(toURL, response)
}

// JSONMultipartForm 上传文件 resp:json
func JSONMultipartForm(toURL string, fields []*MultipartFormField, response interface{}) error {
	return defaultClient().JSONMultipartForm(toURL, fields, response)
}

// Get method:get resp:json
func Get(toURL string, v ...interface{}) *Request {
	return defaultClient().Get(toURL, v...)
}

// Post method:get resp:json
func Post(toURL string, v ...interface{}) *Request {
	return defaultClient().Post(toURL, v...)
}
