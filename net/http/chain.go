package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

// Get 使用get方法
func (client *Client) Get(toURL string, v ...interface{}) *Request {
	u, err := url.Parse(toURL)
	if err != nil {
		return &Request{
			Error: err,
		}
	}
	req := &Request{
		Client: client,
		URL:    u,
		method: "GET",
	}
	for _, vv := range v {
		req.AddRawQuery(vv)
	}
	return req
}

// Post 使用post方法
func (client *Client) Post(toURL string, v ...interface{}) *Request {
	u, err := url.Parse(toURL)
	if err != nil {
		return &Request{
			Error: err,
		}
	}
	req := &Request{
		Client: client,
		URL:    u,
		method: "POST",
	}
	if len(v) > 0 {
		var vv []byte
		switch v[0].(type) {
		case []byte:
			vv = v[0].([]byte)
		case string:
			vv = []byte(v[0].(string))
		default:
			var err error
			vv, err = json.Marshal(v[0])
			if err != nil {
				return &Request{
					Error: err,
				}
			}
		}
		req.postData = vv
	}
	return req
}

// Request http.Request
type Request struct {
	*Client
	*url.URL
	*http.Request
	method   string
	postData []byte
	Error    error
}

// AddRawQuery 设置get form请求
func (req *Request) AddRawQuery(v interface{}) *Request {
	if req.Error != nil {
		return req
	}
	var vv string
	switch v.(type) {
	case string:
		vv = "&" + v.(string)
	default:
		var err error
		vv, err = ToRawQuery(v)
		if err != nil {
			req.Error = err
			return req
		}
	}
	req.URL.RawQuery += "&" + vv
	return req
}

func (req *Request) httpRequest() *Request {
	if req.Error != nil {
		return req
	}
	if req.Request != nil {
		return req
	}
	httpReq, err := http.NewRequest(req.method, req.URL.String(), strings.NewReader(string(req.postData)))
	if err != nil {
		req.Error = err
		return req
	}
	req.Request = httpReq
	return req
}

// JSONRequest json请求
func (req *Request) JSONRequest() *Request {
	req.httpRequest()
	req.SetHeader("Content-Type", "application/json;charset=utf-8")
	return req
}

// FormRequest json请求
func (req *Request) FormRequest() *Request {
	req.httpRequest()
	req.SetHeader("Content-Type", "application/x-www-form-urlencode;charset=utf-8")
	return req
}

// SetHeader 设置请求头
func (req *Request) SetHeader(key, val string) *Request {
	req.httpRequest()
	if req.Error != nil {
		return req
	}
	req.Request.Header.Set(key, val)
	return req
}

// AddHeader 添加请求头
func (req *Request) AddHeader(key, val string) *Request {
	req.httpRequest()
	if req.Error != nil {
		return req
	}
	req.Request.Header.Add(key, val)
	return req
}

// Response 获取响应
func (req *Request) Response() (*http.Response, error) {
	req.httpRequest()
	if req.Error != nil {
		return nil, req.Error
	}
	return req.Client.httpClient().Do(req.Request)
}

// DataResponse 获取响应二进制响应
func (req *Request) DataResponse() ([]byte, error) {
	resp, err := req.Response()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}

// StringResponse 获取响应字符串响应
func (req *Request) StringResponse() (string, error) {
	data, err := req.DataResponse()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// JSONResponse 获取响应JSON响应
func (req *Request) JSONResponse(v interface{}) error {
	data, err := req.DataResponse()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
