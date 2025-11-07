package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	// "suntech.com.vn/skylib/skylog.git/skylog"
)

const DefaultRequestTimeout = 15 * time.Second

type RequestUtil struct {
	Method             string
	Url                string
	Data               map[string]interface{}
	Timeout            time.Duration
	Authorization      string
	Headers            map[string]string
	Transport          *http.Transport
	InsecureSkipVerify int // -1: SkipVerify; 0: Keep default; 1: NotSkipVerify
}

func NewRequest(method, url string, data map[string]interface{}) *RequestUtil {
	r := RequestUtil{}
	r.Method = method
	r.Url = url
	r.Data = data
	r.Timeout = DefaultRequestTimeout
	r.Authorization = ""
	r.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	r.Transport = http.DefaultTransport.(*http.Transport).Clone()
	r.InsecureSkipVerify = 0
	return &r
}

func NewRequestFromJson(jsonStr string) (*RequestUtil, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}
	return NewRequestFromMap(data)
}

func NewRequestFromMap(data map[string]interface{}) (*RequestUtil, error) {
	// Timeout config
	timeoutStr := GetStrValByKey(data, "timeout", "60")
	reqTimeout := time.Duration(ToInt(timeoutStr, 60)) * time.Second

	// Headers config
	headers := map[string]string{}
	if headersRaw, ok := data["headers"]; ok {
		if headersRawMap, ok := headersRaw.(map[string]interface{}); ok {
			for k, v := range headersRawMap {
				headers[k] = ToStr(v, "")
			}
		}
		if headersString, ok := headersRaw.(string); ok {
			var headersData map[string]interface{}
			if err := json.Unmarshal([]byte(headersString), &headersData); err == nil {
				for k, v := range headersData {
					headers[k] = ToStr(v, "")
				}
			}
		}
	}

	// Headers config
	reqData := map[string]interface{}{}
	if reqDataRaw, ok := data["data"]; ok {
		if reqDataRawMap, ok := reqDataRaw.(map[string]interface{}); ok {
			for k, v := range reqDataRawMap {
				reqData[k] = v
			}
		}
		if reqDataString, ok := reqDataRaw.(string); ok {
			var reqDataMap map[string]interface{}
			if err := json.Unmarshal([]byte(reqDataString), &reqDataMap); err == nil {
				for k, v := range reqDataMap {
					reqData[k] = v
				}
			}
		}
	}

	// Do convert
	return &RequestUtil{
		Method:        GetStrValByKey(data, "method", "POST"),
		Url:           GetStrValByKey(data, "url", ""),
		Data:          reqData,
		Timeout:       reqTimeout,
		Authorization: GetStrValByKey(data, "authorization", ""),
		Headers:       headers,
	}, nil
}

func (r *RequestUtil) AddAccessToken(token string) *RequestUtil {
	token = strings.TrimSpace(token)
	if token != "" {
		r.Authorization = token
		if !strings.HasPrefix(token, "Bearer ") {
			token = fmt.Sprintf("Bearer %v", token)
		}
		r.Headers["Authorization"] = token
	}
	return r
}

func (r *RequestUtil) ToRequest(body io.Reader) (*http.Request, error) {
	// create request
	req, err := http.NewRequest(r.Method, r.Url, body)
	if err != nil {
		// skylog.Errorf("SendRawFromRequest create request error: %v", err)
		fmt.Println("SendRawFromRequest create request error: %v", err)
		return nil, err
	}

	// check authorization
	r.Authorization = strings.TrimSpace(r.Authorization)
	if r.Authorization != "" {
		if _, ok := r.Headers["Authorization"]; !ok {
			if !strings.HasPrefix(r.Authorization, "Bearer ") {
				r.Headers["Authorization"] = fmt.Sprintf("Bearer %v", r.Authorization)
			}
		}
	}

	// setup header request
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}
	return req, nil
}

func (r *RequestUtil) ToClient() *http.Client {
	// create http client
	if r.InsecureSkipVerify < 0 {
		r.Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else if r.InsecureSkipVerify > 0 {
		r.Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	}
	return &http.Client{
		Transport: r.Transport,
		Timeout:   r.Timeout,
	}
}

func (r *RequestUtil) SendRaw() (*http.Response, error) {
	// data
	var data io.Reader = nil
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		bodyData, err := json.Marshal(r.Data)
		if err != nil {
			// skylog.Errorf("RequestUtil Send create body data error: %v | data: %v", err, r.Data)
			fmt.Println("RequestUtil Send create body data error: %v | data: %v", err, r.Data)
			return nil, err
		}
		data = bytes.NewReader(bodyData)
	}

	// create request
	req, err := r.ToRequest(data)
	if err != nil {
		return nil, err
	}
	// create http client
	if r.InsecureSkipVerify < 0 {
		r.Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else if r.InsecureSkipVerify > 0 {
		r.Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	}
	client := &http.Client{
		Transport: r.Transport,
		Timeout:   r.Timeout,
	}
	// send request
	res, err := client.Do(req)
	if res != nil && res.StatusCode != http.StatusOK {
		bodyData, _ := json.Marshal(r.Data)
		// skylog.Infof("SendRequest [%v]%v -> %v\nData: %v\n", r.Method, r.Url, res.Status, string(bodyData))
		fmt.Println("SendRequest [%v]%v -> %v\nData: %v\n", r.Method, r.Url, res.Status, string(bodyData))
	}

	return res, err
}

func (r *RequestUtil) Send() ([]byte, bool, error) {
	// send request
	res, err := r.SendRaw()
	if err != nil {
		return nil, false, err
	}
	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, false, err
	}

	if res.StatusCode == http.StatusOK {
		return resData, true, nil
	}
	return resData, false, nil
}
