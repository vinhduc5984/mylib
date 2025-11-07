package utils

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"mime/multipart"
// 	"net/http"
// 	"net/url"
// 	"strings"

// 	"suntech.com.vn/skylib/skylog.git/skylog"
// )

// type ReadDocInfo struct {
// 	CompanyId   int64  `json:"companyId"`
// 	BranchId    int64  `json:"branchId"`
// 	Iuid        string `json:"iuid"`
// 	ServiceCode string `json:"serviceCode"`
// 	ScreenCode  string `json:"screenCode"`
// 	FeatureCode string `json:"featureCode"`
// 	FullPath    string `json:"fullPath"`
// 	Mode        string `json:"mode"`
// }

// type WriteDocInfo struct {
// 	CompanyId   int64  `json:"companyId"`
// 	BranchId    int64  `json:"branchId"`
// 	Iuid        string `json:"iuid"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Note        string `json:"note"`
// 	ServiceCode string `json:"serviceCode"`
// 	ScreenCode  string `json:"screenCode"`
// 	FeatureCode string `json:"featureCode"`
// 	ItemType    string `json:"itemType"`
// 	ItemCode    string `json:"itemCode"`
// 	ItemId      int64  `json:"itemId"`
// 	ItemDate    int64  `json:"itemDate"`
// 	PartnerCode string `json:"partnerCode"`
// 	PartnerName string `json:"partnerName"`
// 	Mode        int32  `json:"mode"`
// 	Md5         string `json:"md5"`
// 	Size        int64  `json:"size"`
// 	Force       bool   `json:"force"`
// }

// type ShareDocInfo struct {
// 	Id    string `json:"id"`
// 	Name  string `json:"name"`
// 	Token string `json:"token"`
// 	Mode  string `json:"mode"`
// 	Url   string `json:"url"`
// }

// func WriteDoc(ctx context.Context, urlOrServiceAddress string, metadata WriteDocInfo, data []byte) ([]byte, error) {
// 	payload := &bytes.Buffer{}
// 	writer := multipart.NewWriter(payload)
// 	fileData, err := writer.CreateFormFile("data", metadata.Name)
// 	if err != nil {
// 		skylog.Errorf("WriteBusinessFile create form file error: %v", err)
// 		return nil, err
// 	}

// 	// convert byte slice to io.Reader
// 	reader := bytes.NewReader(data)
// 	_, err = io.Copy(fileData, reader)
// 	if err != nil {
// 		skylog.Errorf("WriteBusinessFile copy byte array error error: %v", err)
// 		return nil, err
// 	}
// 	_ = writer.WriteField("companyId", ToStr(metadata.CompanyId))
// 	_ = writer.WriteField("branchId", ToStr(metadata.BranchId))
// 	_ = writer.WriteField("iuid", metadata.Iuid)
// 	_ = writer.WriteField("name", metadata.Name)
// 	_ = writer.WriteField("description", metadata.Description)
// 	_ = writer.WriteField("note", metadata.Note)
// 	_ = writer.WriteField("serviceCode", metadata.ServiceCode)
// 	_ = writer.WriteField("screenCode", metadata.ScreenCode)
// 	_ = writer.WriteField("featureCode", metadata.FeatureCode)
// 	_ = writer.WriteField("itemType", metadata.ItemType)
// 	_ = writer.WriteField("itemCode", metadata.ItemCode)
// 	_ = writer.WriteField("itemId", ToStr(metadata.ItemId))
// 	_ = writer.WriteField("itemDate", ToStr(metadata.ItemDate))
// 	_ = writer.WriteField("partnerCode", metadata.PartnerCode)
// 	_ = writer.WriteField("partnerName", metadata.PartnerName)
// 	_ = writer.WriteField("mode", ToStr(metadata.Mode))
// 	_ = writer.WriteField("md5", metadata.Md5)
// 	_ = writer.WriteField("size", ToStr(metadata.Size))
// 	_ = writer.WriteField("force", ToStr(metadata.Force))

// 	err = writer.Close()
// 	if err != nil {
// 		skylog.Errorf("WriteBusinessFile close write error: %v", err)
// 		return nil, err
// 	}

// 	urlPath := "/doc/file/v1/upload"
// 	return RestUploadFile(urlOrServiceAddress, urlPath, payload, writer.FormDataContentType(), ctx)
// }

// func ReadDocWithUrl(ctx context.Context, fullUrl string) ([]byte, string, string, string, error) {
// 	return RestDownloadFile(fullUrl, "", ctx)
// }

// func ReadDoc(ctx context.Context, urlOrServiceAddress string, metadata ReadDocInfo) ([]byte, string, string, string, error) {
// 	urlPath := ""
// 	fullPath := strings.TrimSpace(metadata.FullPath)
// 	if len(fullPath) > 0 {
// 		urlPath = BuildDocUrlWithFullPath(fullPath)
// 	} else {
// 		urlPath = BuildDocUrl(metadata.Iuid, metadata.CompanyId, metadata.BranchId, metadata.ServiceCode, metadata.ScreenCode, metadata.FeatureCode)
// 	}
// 	return RestDownloadFile(urlOrServiceAddress, urlPath, ctx)
// }

// func BuildDocUrlWithData(params url.Values, opts ...interface{}) string {
// 	accessMethod := "download"
// 	if len(opts) > 0 {
// 		tmp := opts[0].(string)
// 		if tmp == "view" || tmp == "preview" {
// 			accessMethod = tmp
// 		}
// 	}
// 	checksumUrl := fmt.Sprintf("/doc/file/v1/%v?%v", accessMethod, params.Encode())
// 	checksum := Sha256(checksumUrl)
// 	params.Add("checksum", checksum)
// 	return fmt.Sprintf("/doc/file/v1/%v?%v", accessMethod, params.Encode())
// }

// func BuildDocUrlWithFullPath(fullPath string, opts ...interface{}) string {
// 	params := url.Values{}
// 	params.Add("fullPath", fullPath)

// 	return BuildDocUrlWithData(params, opts...)
// }

// func BuildDocUrl(iuid string, companyId, branchId int64, serviceCode, screenCode, featureCode string, opts ...interface{}) string {
// 	params := url.Values{}
// 	params.Add("iuid", iuid)
// 	params.Add("companyId", ToStr(companyId, "0"))
// 	params.Add("branchId", ToStr(branchId, "0"))
// 	params.Add("service", serviceCode)
// 	params.Add("screen", screenCode)
// 	params.Add("feature", featureCode)

// 	return BuildDocUrlWithData(params, opts...)
// }

// func BuildDocUrlWithToken(documentId int64, token string, opts ...interface{}) string {
// 	params := url.Values{}
// 	params.Add("id", ToStr(documentId, ""))
// 	params.Add("token", token)

// 	return BuildDocUrlWithData(params, opts...)
// }

// func ShareDoc(ctx context.Context, urlOrServiceAddress string, data ReadDocInfo) (map[string]interface{}, error) {
// 	return ShareDocWithContext(urlOrServiceAddress, data, ctx)
// }

// func ShareDocWithAccessToken(urlOrServiceAddress string, data ReadDocInfo, accessToken string) (map[string]interface{}, error) {
// 	var returnVal map[string]interface{}
// 	return returnVal, nil
// }

// func ShareDocWithContext(urlOrServiceAddress string, data ReadDocInfo, ctx context.Context) (map[string]interface{}, error) {
// 	params := map[string]interface{}{
// 		"companyId":   data.CompanyId,
// 		"branchId":    data.BranchId,
// 		"serviceCode": data.ServiceCode,
// 		"screenCode":  data.ScreenCode,
// 		"featureCode": data.FeatureCode,
// 		"iuid":        data.Iuid,
// 		"mode":        data.Mode,
// 	}

// 	res, err := RestPostWithContext(urlOrServiceAddress, "/doc/file/v1/get-link", params, ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var returnVal map[string]interface{}
// 	if err := json.Unmarshal([]byte(res), &returnVal); err != nil {
// 		return nil, err
// 	}
// 	return returnVal, nil
// }

// func ResponseBusinessFile(w http.ResponseWriter, docServiceAddress string, ctx context.Context, metadata ReadDocInfo) error {
// 	urlPath := ""
// 	fullPath := strings.TrimSpace(metadata.FullPath)
// 	if len(fullPath) > 0 {
// 		urlPath = BuildDocUrlWithFullPath(fullPath)
// 	} else {
// 		urlPath = BuildDocUrl(metadata.Iuid, metadata.CompanyId, metadata.BranchId, metadata.ServiceCode, metadata.ScreenCode, metadata.FeatureCode)
// 	}

// 	accessToken, _, _ := GetLoginAccessToken(ctx)
// 	resp, err := SendRawFromRequest(http.MethodGet, docServiceAddress, urlPath, accessToken, nil, ctx)
// 	if err != nil {
// 		ResponseInternalError(w, err.Error())
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	return ForwardResponse(resp, w)
// }

// func DeleteDoc(ctx context.Context, urlOrServiceAddress string, data ReadDocInfo) error {
// 	return DeleteDocWithAccessToken(urlOrServiceAddress, data, ctx)
// }

// func DeleteDocWithAccessToken(urlOrServiceAddress string, data ReadDocInfo, ctx context.Context) error {
// 	iuid := []string{data.Iuid}
// 	params := map[string]interface{}{
// 		"companyId":   data.CompanyId,
// 		"branchId":    data.BranchId,
// 		"serviceCode": data.ServiceCode,
// 		"screenCode":  data.ScreenCode,
// 		"featureCode": data.FeatureCode,
// 		"iuid":        iuid,
// 	}

// 	_, err := RestPostWithContext(urlOrServiceAddress, "/doc/file/v1/delete-business", params, ctx)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // func ReplicateFile(ctx context.Context, urlOrServiceAddr string) error {
// // 	accessToken, _, _ := GetLoginAccessToken(ctx)
// // 	res, err := RestPost(urlOrServiceAddr, "/core/notification/v1/send", map[string]interface{}{}, accessToken)
// // 	if err != nil {
// // 		return err
// // 	}

// // 	return nil
// // }
