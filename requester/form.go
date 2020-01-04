package requester

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
)

// FormFile 文件表单数据
type FormFile struct {
	FieldName string `json:"field_name,omitempty"`
	FilePath  string `json:"file_path,omitempty"`
	FileName  string `json:"file_name,omitempty"`
}

// 注意 如果key重复 会追加成列表
func buildFormData(fields map[string][]string, files []*FormFile) (string, io.Reader, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for key, val := range fields {
		for _, vv := range val {
			if err := bodyWriter.WriteField(key, vv); err != nil {
				return "", nil, err
			}
		}
	}

	for _, file := range files {
		fileWriter, err := bodyWriter.CreateFormFile(file.FieldName, file.FileName)
		if err != nil {
			return "", nil, err
		}

		fh, err := os.Open(file.FilePath)
		if err != nil {
			return "", nil, err
		}
		defer fh.Close()

		if _, err = io.Copy(fileWriter, fh); err != nil {
			return "", nil, err
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	return contentType, bodyBuf, nil
}
