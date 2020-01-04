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

// JSON method:post req:json resp:json
func (client *Client) JSON(toURL string, request interface{}, response interface{}) error {
	postBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", toURL, strings.NewReader(string(postBody)))
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

// GetJSON method:get resp:json
func (client *Client) GetJSON(toURL string, response interface{}) error {
	resp, err := client.httpClient().Get(toURL)
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

// PostJSON method:get resp:json
func (client *Client) PostJSON(toURL string, response interface{}) error {
	req, err := http.NewRequest("POST", toURL, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-toURLencoded;charset=utf-8")

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

// JSONMultipartForm 上传文件 resp:json
func (client *Client) JSONMultipartForm(toURL string, fields []*MultipartFormField, response interface{}) error {
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

	resp, err := client.httpClient().Post(toURL, contentType, bodyBuf)
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
