package http

import (
	"bytes"
	"encoding/json"
	"github.com/waittttting/cRPC-common/cerr"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func Post(url string, params map[string]string) (interface{}, error) {

	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for k, v := range params {
		err := writer.WriteField(k, v)
		if err != nil {
			return nil, err
		}
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp CResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code == BusinessOk {
		return resp.Data, nil
	} else {
		dataMap := resp.Data.(map[string]interface{})
		return nil, cerr.NewError(int64(dataMap["ErrCode"].(float64)), dataMap["ErrMsg"].(string))
	}
}
