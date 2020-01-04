package requester

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testResponse struct {
	RawQuery string              `json:"rawquery,omitempty"`
	Files    map[string]int64    `json:"files,omitempty"`
	Form     map[string][]string `json:"form,omitempty"`
	Data     []byte              `json:"data,omitempty"`
	Headers  map[string][]string `json:"headers,omitempty"`
}

func testRawServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		resp := &testResponse{
			RawQuery: r.URL.RawQuery,
			Data:     data,
			Headers:  r.Header,
		}
		respData, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		w.Write(respData)
	}))
	return ts
}

func testFormServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := r.ParseMultipartForm(2048)
		if err != nil {
			panic(err)
		}

		resp := &testResponse{
			RawQuery: r.URL.RawQuery,
			Form:     r.MultipartForm.Value,
			Headers:  r.Header,
			Files:    map[string]int64{},
		}
		for key, val := range r.MultipartForm.File {
			resp.Files[key] = val[0].Size
		}

		respData, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		w.Write(respData)
	}))
	return ts
}

func TestRequester(t *testing.T) {
	testKey1 := "abc"
	testKey2 := "def"
	testVal1 := "123"
	testVal2 := "456"

	rawQuery1 := testKey1 + "=" + testVal1
	rawQuery2 := testKey2 + "=" + testVal2
	reqForm := map[string]string{testKey1: testVal1}
	reqData, _ := json.Marshal(reqForm)
	reqFile := &FormFile{
		FieldName: testKey2,
		FileName:  testVal2,
		FilePath:  "./testdata/abc.txt",
	}

	t.Run("Get", func(t *testing.T) {
		ts := testRawServer()
		defer ts.Close()

		resp := &testResponse{}
		err := New("GET", ts.URL+"?"+rawQuery1).AddRawQuery(rawQuery2).ReadJSON(resp)
		if err != nil {
			t.Error(err)
		}
		if resp.RawQuery != rawQuery1+"&"+rawQuery2 {
			t.Errorf("RawQuery error")
		}
	})
	t.Run("PostRaw", func(t *testing.T) {
		ts := testRawServer()
		defer ts.Close()

		resp := &testResponse{}
		err := New("POST", ts.URL+"?"+rawQuery1).SetRawData(reqData).ReadJSON(resp)
		if err != nil {
			t.Error(err)
		}
		if resp.RawQuery != rawQuery1 {
			t.Errorf("RawQuery error")
		}
		if string(resp.Data) != string(reqData) {
			t.Errorf("Data error")
		}
	})
	t.Run("PostForm", func(t *testing.T) {
		ts := testFormServer()
		defer ts.Close()

		resp := &testResponse{}
		err := New("POST", ts.URL+"?"+rawQuery1).SetFormFields(reqForm).ReadJSON(resp)
		if err != nil {
			t.Error(err)
		}
		if resp.RawQuery != rawQuery1 {
			t.Errorf("RawQuery error")
		}
		v, ok := resp.Form[testKey1]
		if !ok {
			t.Errorf("Form error")
		}
		if v[0] != testVal1 {
			t.Errorf("Form Val error")
		}
	})
	t.Run("PostFile", func(t *testing.T) {
		ts := testFormServer()
		defer ts.Close()

		resp := &testResponse{}
		err := New("POST", ts.URL+"?"+rawQuery1).SetFormFields(reqForm).AddFormFile(reqFile).ReadJSON(resp)
		if err != nil {
			t.Error(err)
		}
		if resp.RawQuery != rawQuery1 {
			t.Errorf("RawQuery error")
		}
		v, ok := resp.Form[testKey1]
		if !ok {
			t.Errorf("Form error")
		}
		if v[0] != testVal1 {
			t.Errorf("Form Val error")
		}
		f, ok := resp.Files[testKey2]
		if !ok {
			t.Errorf("File error")
		}
		if f == 0 {
			t.Errorf("File Size error")
		}
	})
}
