package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"
	// "suntech.com.vn/skylib/skylog.git/skylog"
)

const (
	BUFFER_SIZE = 1024 * 4
)

func BuildServiceUrl(urlOrServiceAddr, path string) (string, error) {
	if strings.HasPrefix(urlOrServiceAddr, "http") {
		return urlOrServiceAddr, nil
	}

	parts := strings.Split(urlOrServiceAddr, ":")
	if len(parts) < 2 {
		// skylog.Errorf("service address is incorrect: %v", urlOrServiceAddr)
		fmt.Println("service address is incorrect: %v", urlOrServiceAddr)
		return "", errors.New("service address is incorrect")
	}
	grpcPort, _ := strconv.Atoi(parts[1])

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%v", path)
	}
	url := fmt.Sprintf("http://%v:%v%v", parts[0], grpcPort+1, path)

	return url, nil
}

func SendRawFromRequest(method, urlOrServiceAddr, path, accessToken string, body map[string]interface{}, callFrom interface{}) (*http.Response, error) {
	// build url
	serviceUrl, err := BuildServiceUrl(urlOrServiceAddr, path)
	if err != nil {
		// skylog.Errorf("SendRawFromRequest build service url error: %v", err)
		fmt.Println("SendRawFromRequest build service url error: %v", err)
		return nil, err
	}

	// data
	var data io.Reader = nil
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		bodyData, err := json.Marshal(body)
		if err != nil {
			// skylog.Errorf("SendRawFromRequest create body data error: %v | data: %v", err, body)
			fmt.Println("SendRawFromRequest create body data error: %v | data: %v", err, body)
			return nil, err
		}
		data = bytes.NewReader(bodyData)
	}
	// create request
	req, err := http.NewRequest(method, serviceUrl, data)
	if err != nil {
		// skylog.Errorf("SendRawFromRequest create request error: %v", err)
		fmt.Println("SendRawFromRequest create request error: %v", err)
		return nil, err
	}

	// set User-Agent
	if callFrom != nil {
		if reflect.TypeOf(callFrom).String() == "*http.Request" {
			fromReq := callFrom.(*http.Request)
			req.Header.Set("User-Agent", fromReq.Header.Get("User-Agent"))

			xRealIP := fromReq.Header.Get("X-Real-IP")
			if len(xRealIP) > 0 {
				req.Header.Set("X-Real-IP", xRealIP)
			}

			xForwardedFor := fromReq.Header.Get("X-Forwarded-For")
			if len(xForwardedFor) > 0 {
				req.Header.Set("X-Forwarded-For", xForwardedFor)
			}

			xForwardedHost := fromReq.Header.Get("X-Forwarded-Host")
			if len(xForwardedHost) > 0 {
				req.Header.Set("X-Forwarded-Host", xForwardedHost)
			} else {
				req.Header.Set("X-Forwarded-Host", fromReq.Host)
			}

			if len(fromReq.RemoteAddr) > 0 {
				req.RemoteAddr = fromReq.RemoteAddr
			}
		} else {
			fromContext := callFrom.(context.Context)
			// Extract metadata from gRPC context
			md, ok := metadata.FromIncomingContext(fromContext)
			if !ok {
				md = metadata.New(nil)
			}
			// Add metadata as HTTP headers
			for key, values := range md {
				// Skip pseudo-headers like ":authority" or other invalid header names
				if key == ":authority" || key == ":method" || key == ":path" || key == ":scheme" {
					continue
				}
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		}
	}

	// set request param
	if len(accessToken) > 0 {
		if !strings.HasPrefix(accessToken, "Bearer ") {
			accessToken = fmt.Sprintf("Bearer %v", accessToken)
		}
		req.Header.Set("Authorization", accessToken)
	}
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}

	// create http client
	client := http.Client{}
	// send request
	res, err := client.Do(req)
	return res, err
}

func SendRawRequest(method, urlOrServiceAddr, path, accessToken string, body map[string]interface{}) (*http.Response, error) {
	// send request
	return SendRawFromRequest(method, urlOrServiceAddr, path, accessToken, body, nil)
}

func SendMultiPartForm(method, urlOrServiceAddr, path string, body io.Reader, contentType string, ctx context.Context) (*http.Response, error) {
	// build url
	serviceUrl, err := BuildServiceUrl(urlOrServiceAddr, path)
	if err != nil {
		// skylog.Errorf("SendMultiPartForm build service url error: %v", err)
		fmt.Println("SendMultiPartForm build service url error: %v", err)
		return nil, err
	}

	// create request
	req, err := http.NewRequest(method, serviceUrl, body)
	if err != nil {
		// skylog.Errorf("SendMultiPartForm create request error: %v", err)
		fmt.Println("SendMultiPartForm create request error: %v", err)
		return nil, err
	}

	// Extract metadata from gRPC context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	// Add metadata as HTTP headers
	for key, values := range md {
		// Skip pseudo-headers like ":authority" or other invalid header names
		if key == ":authority" || key == ":method" || key == ":path" || key == ":scheme" {
			continue
		}
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	accessToken, _, _ := GetLoginAccessToken(ctx)
	// set request param
	if len(accessToken) > 0 {
		if !strings.HasPrefix(accessToken, "Bearer ") {
			accessToken = fmt.Sprintf("Bearer %v", accessToken)
		}
		req.Header.Set("Authorization", accessToken)
	}
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req.Header.Set("Content-Type", contentType)
	}

	// create http client
	client := http.Client{}
	// send request
	return client.Do(req)
}

func ForwardRequest(method, urlOrServiceAddr, path, accessToken string, body map[string]interface{}, fromReqOrContext interface{}) ([]byte, error) {
	// send request
	resp, err := SendRawFromRequest(method, urlOrServiceAddr, path, accessToken, body, fromReqOrContext)
	if err != nil {
		// skylog.Errorf("SendRequest send request error 1: %v", err)
		fmt.Println("SendRequest send request error 1: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	resData, err := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return resData, nil
	}
	// skylog.Errorf("SendRequest send request error 2: %v", string(resData))
	fmt.Println("SendRequest send request error 2: %v", string(resData))
	return resData, CallRESTAPIError
}

func SendRequest(method, urlOrServiceAddr, path, accessToken string, body map[string]interface{}) ([]byte, error) {
	return ForwardRequest(method, urlOrServiceAddr, path, accessToken, body, nil)
}

func ForwardGet(urlOrServiceAddr, path string, accessToken string, fromReq *http.Request) ([]byte, error) {
	return ForwardRequest(http.MethodGet, urlOrServiceAddr, path, accessToken, nil, fromReq)
}

func RestGet(urlOrServiceAddr, path string, accessToken string) ([]byte, error) {
	return ForwardRequest(http.MethodGet, urlOrServiceAddr, path, accessToken, nil, nil)
}

func ForwardPost(urlOrServiceAddr, path string, values map[string]interface{}, accessToken string, fromReq *http.Request) ([]byte, error) {
	return ForwardRequest(http.MethodPost, urlOrServiceAddr, path, accessToken, values, fromReq)
}

func RestPost(urlOrServiceAddr, path string, values map[string]interface{}, accessToken string) ([]byte, error) {
	return ForwardRequest(http.MethodPost, urlOrServiceAddr, path, accessToken, values, nil)
}

func RestPostWithContext(urlOrServiceAddr, path string, values map[string]interface{}, ctx context.Context) ([]byte, error) {
	accessToken, _, _ := GetLoginAccessToken(ctx)
	return ForwardRequest(http.MethodPost, urlOrServiceAddr, path, accessToken, values, ctx)
}

func RestDownloadFile(urlOrServiceAddr, path string, ctx context.Context) ([]byte, string, string, string, error) {
	accessToken, _, _ := GetLoginAccessToken(ctx)
	// send request
	resp, err := SendRawFromRequest(http.MethodGet, urlOrServiceAddr, path, accessToken, nil, ctx)
	if err != nil {
		// skylog.Errorf("RestDownloadFile send request error: %v", err)
		fmt.Println("RestDownloadFile send request error: %v", err)
		return nil, "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Get the file name from the Content-Disposition header
		fileName := resp.Header.Get("Content-Disposition")
		fileName = strings.TrimSpace(strings.Split(fileName, ";")[1])
		fileName = strings.TrimSpace(strings.Split(fileName, "=")[1])
		fileName = strings.Replace(fileName, "\"", "", -1)

		// Get the content type from the Content-Type header
		contentType := resp.Header.Get("Content-Type")

		// Get the date from the Last-Modified header
		dateStr := resp.Header.Get("Last-Modified")
		if dateStr != "" {
			date, err := time.Parse(http.TimeFormat, dateStr)
			if err != nil {
				// skylog.Errorf("RestDownloadFile parse time error; time=%v; err:%v", dateStr, err)
				fmt.Println("RestDownloadFile parse time error; time=%v; err:%v", dateStr, err)
				dateStr = ""
			} else {
				dateStr = ToStr(date.UTC().UnixMilli())
			}
		}

		// Read the response body into a byte array
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			// skylog.Errorf("RestDownloadFile read body error: %v", err)
			fmt.Println("RestDownloadFile read body error: %v", err)
			return nil, fileName, contentType, dateStr, err
		}

		return bytes, fileName, contentType, dateStr, nil
	} else {
		// skylog.Errorf("RestDownloadFile call request error: %v", resp.StatusCode)
		fmt.Println("RestDownloadFile call request error: %v", resp.StatusCode)
		// Read the response body into a byte array
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			// skylog.Errorf("RestDownloadFile read body error: %v", err)
			fmt.Println("RestDownloadFile read body error: %v", err)
			return nil, "", "", "", err
		}

		return bytes, "", "", "", errors.New(string(bytes))
	}
}

func RestUploadFile(urlOrServiceAddr, path string, body io.Reader, contentType string, ctx context.Context) ([]byte, error) {
	// send request
	resp, err := SendMultiPartForm(http.MethodPost, urlOrServiceAddr, path, body, contentType, ctx)
	if err != nil {
		// skylog.Errorf("RestUploadFile send multi part file error: %v", err)
		fmt.Println("RestUploadFile send multi part file error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	resData, err := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return resData, nil
	}
	// skylog.Errorf("RestUploadFile send multi part file error: %v", string(resData))
	fmt.Println("RestUploadFile send multi part file error: %v", string(resData))
	return resData, CallRESTAPIError
}

func IsValidChecksum(r *http.Request) bool {
	// get client checksum from url
	checksum := r.URL.Query().Get("checksum")
	if len(checksum) > 0 {
		// remove checksum from url
		queryParams := r.URL.Query()
		queryParams.Del("checksum")

		// build checksum url
		checksumUrl := fmt.Sprintf("%v?%v", r.URL.Path, queryParams.Encode())
		serverChecksum := Sha256(checksumUrl)

		if checksum == serverChecksum {
			return true
		}
		// skylog.Infof("IsValidChecksum: checksum = %v <> %v = serverChecksum (%v)", checksum, serverChecksum, checksumUrl)
		fmt.Println("IsValidChecksum: checksum = %v <> %v = serverChecksum (%v)", checksum, serverChecksum, checksumUrl)
	}
	return false
}

func IsValidChecksumWithContext(ctx context.Context) bool {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if url, hasUrl := md["pattern"]; hasUrl && len(url) > 0 {
			return IsValidChecksumWithUrl(url[0])
		} else {
			fmt.Printf("md: %v\n", md)
		}
	}
	return false
}

func IsValidChecksumWithUrl(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	// get client checksum from url
	checksum := u.Query().Get("checksum")
	if len(checksum) > 0 {
		// remove checksum from url
		queryParams := u.Query()
		queryParams.Del("checksum")

		// build checksum url
		checksumUrl := fmt.Sprintf("%v?%v", u.Path, queryParams.Encode())
		serverChecksum := Sha256(checksumUrl)

		if checksum == serverChecksum {
			return true
		}
		// skylog.Infof("IsValidChecksumWithUrl: client checksum = %v <> %v = serverChecksum (%v)", checksum, serverChecksum, checksumUrl)
		fmt.Println("IsValidChecksumWithUrl: client checksum = %v <> %v = serverChecksum (%v)", checksum, serverChecksum, checksumUrl)
	}
	return false
}

func BuildQrCodeUrlWithId(ctx context.Context, id int64) (*string, error) {
	var scheme, host []string
	var hasValue bool
	// get metadata from context
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// get schema value
		if scheme, hasValue = md["x-forwarded-scheme"]; !hasValue || len(scheme) == 0 {
			scheme = []string{"http"}
		}
		// get host value
		if host, hasValue = md["x-forwarded-host"]; !hasValue || len(host) == 0 {
			return nil, Error500WithMessage("SYS.MSG.MISSING_HOST_FROM_REQUEST_HEADER")
		}
	} else {
		return nil, Error500WithMessage("SYS.MSG.MISSING_REQUEST_INFO_FROM_HEADER")
	}

	// build url
	returnVal := fmt.Sprintf("%v://%v/orbit/%v", scheme[0], host[0], EncryptId(id))
	return &returnVal, nil
}

func BuildQrCodeUrlFromRequest(req *http.Request, id int64) (*string, error) {
	if req == nil {
		return nil, Error500WithMessage("SYS.MSG.MISSING_REQUEST")
	}
	// check host
	host := req.Header.Get("X-Forwarded-Host")
	if strings.TrimSpace(host) == "" {
		// check origin
		origin := req.Header.Get("Origin")
		if strings.TrimSpace(origin) == "" {
			return nil, Error500WithMessage("SYS.MSG.MISSING_HOST_FROM_REQUEST_HEADER")
		} else {
			returnVal := fmt.Sprintf("%v/orbit/%v", origin, EncryptId(id))
			return &returnVal, nil
		}
	}
	// check scheme
	scheme := req.Header.Get("X-Forwarded-Scheme")
	if strings.TrimSpace(scheme) == "" {
		scheme = "http"
	}

	// build url
	returnVal := fmt.Sprintf("%v://%v/orbit/%v", scheme, host, EncryptId(id))
	return &returnVal, nil
}

func ForwardResponse(src *http.Response, dest http.ResponseWriter) error {
	// Create a buffer to store the file data.
	buffer := make([]byte, BUFFER_SIZE)

	dest.Header().Set("Content-Type", src.Header.Get("Content-Type"))
	dest.Header().Set("Content-Length", src.Header.Get("Content-Length"))
	dest.Header().Set("Last-Modified", src.Header.Get("Last-Modified"))

	if _, err := io.CopyBuffer(dest, src.Body, buffer); err != nil {
		ResponseInternalError(dest, err.Error())
		return err
	}
	return nil
}
