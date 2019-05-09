package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// JSON post 方式获取 json返回
func JSON(url string, request interface{}, response interface{}) error {
	postBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(postBody)))
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		return err
	}
	return nil
}

// GetJSON 用get方式获取json返回
func GetJSON(url string, response interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		return err
	}
	return nil
}
