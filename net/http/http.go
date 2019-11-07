// Package http 网络请求方法
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client Client
type Client struct {
	Transport     http.RoundTripper
	CheckRedirect func(req *http.Request, via []*http.Request) error
	Jar           http.CookieJar
	Timeout       time.Duration
}

func (client *Client) httpClient() *http.Client {
	return &http.Client{
		Transport:     client.Transport,
		CheckRedirect: client.CheckRedirect,
		Jar:           client.Jar,
		Timeout:       client.Timeout,
	}
}

// JSON post 方式获取 json返回
func (client *Client) JSON(url string, request interface{}, response interface{}) error {
	postBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(postBody)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", resp.Status)
	}

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

// GetJSON 用get方式获取json返回
func (client *Client) GetJSON(url string, response interface{}) error {
	resp, err := client.httpClient().Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", resp.Status)
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

// PostJSON 用post方式获取json返回
func (client *Client) PostJSON(url string, response interface{}) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := client.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", resp.Status)
	}

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

// MultipartFormField 文件或其他表单数据
type MultipartFormField struct {
	FieldName  string
	FieldValue string
	FilePath   string
	FileName   string
}

// JSONMultipartForm 上传文件或其他表单数据
func (client *Client) JSONMultipartForm(url string, fields []*MultipartFormField, response interface{}) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for _, field := range fields {
		if field.FilePath == "" {
			partWriter, err := bodyWriter.CreateFormField(field.FieldName)
			if err != nil {
				return err
			}
			valueReader := bytes.NewReader([]byte(field.FieldValue))
			if _, err = io.Copy(partWriter, valueReader); err != nil {
				return err
			}
		} else {
			fileWriter, err := bodyWriter.CreateFormFile(field.FieldName, field.FileName)
			if err != nil {
				return err
			}

			fh, err := os.Open(field.FilePath)
			if err != nil {
				return err
			}
			defer fh.Close()

			if _, err = io.Copy(fileWriter, fh); err != nil {
				return err
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := client.httpClient().Post(url, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", resp.Status)
	}

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
