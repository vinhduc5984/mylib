package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
	// "suntech.com.vn/skylib/skylog.git/skylog"
)

const (
	DiffHourNil = float64(-99)
)

type LoginInfo struct {
	UserId       int64
	Username     string
	CompanyId    int64
	BranchId     int64
	DepartmentId int64
	PartnerId    int64
	PartnerCode  string
	Ip           string
	PartnerName  string
	DiffHour     float64
	DeviceId     int64
	AccountType  int32
}

// GetAccessToken function return full accessToken from context and jwt manager
func GetAccessToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return "", Unauthenticated
	}

	auths := md["authorization"]
	if len(auths) == 0 || len(auths[0]) == 0 {
		return "", Unauthenticated
	}

	return auths[0], nil
}

// GetLoginAccessToken function return accessToken, suffix token from context and jwt manager
func GetLoginAccessToken(ctx context.Context) (string, string, error) {
	fullAccessToken, err := GetAccessToken(ctx)
	if err != nil {
		return "", "", err
	}

	prefix := "Bearer "
	splitter := "|||"
	suffix := ""
	accessToken := strings.Replace(fullAccessToken, prefix, "", 1)
	index := strings.Index(accessToken, splitter)
	if index >= 0 {
		suffix = accessToken[len(splitter)+index:]
		accessToken = accessToken[:index]
	}
	return accessToken, suffix, nil
}

func getTokenSuffix(accessToken string) (string, string, error) {
	splitter := "|||"
	suffix := ""
	index := strings.Index(accessToken, splitter)
	if index >= 0 {
		suffix = accessToken[len(splitter)+index:]
		accessToken = accessToken[:index]
	}
	return accessToken, suffix, nil
}

// GetUserClaims function return user id from context and jwt manager
func GetUserClaims(ctx context.Context) (*jwt.MapClaims, error) {
	accessToken, _, err := GetLoginAccessToken(ctx)
	if err != nil {
		// skylog.Info("error in GetLoginAccessToken")
		fmt.Println("error in GetLoginAccessToken")
		return nil, err
	}
	// skylog.Info("goto GetUserClaimsFromToken")
	fmt.Println("goto GetUserClaimsFromToken")
	return GetUserClaimsFromToken(accessToken)
}

// ToAccount convert from user claims to Account
func ToAccount(claims *jwt.MapClaims) Account {
	userID, _ := ToInt64((*claims)["userId"])
	partnerId, _ := ToInt64((*claims)["partnerId"])
	partnerCode, _ := ToString((*claims)["partnerCode"])
	partnerName, _ := ToString((*claims)["partnerName"])
	ip, _ := ToString((*claims)["ip"])
	diffHour := ToFloat64WithDefault((*claims)["diffHour"], DiffHourNil)
	username, _ := ToString((*claims)["username"])
	fullName, _ := ToString((*claims)["fullName"])
	deviceId, _ := ToInt64((*claims)["deviceId"])
	accountType, _ := ToInt32((*claims)["accountType"])

	return Account{
		Id:          userID,
		Username:    &username,
		PartnerId:   partnerId,
		PartnerCode: partnerCode,
		PartnerName: partnerName,
		Ip:          ip,
		DiffHour:    diffHour,
		FullName:    fullName,
		DeviceId:    deviceId,
		AccountType: accountType,
	}
}

// GetAccountInfo function return Account from context and jwt manager
func GetAccountInfo(ctx context.Context) (Account, error) {
	userClaims, err := GetUserClaims(ctx)

	if err != nil {
		return Account{}, err
	}

	return ToAccount(userClaims), nil
}

// GetAccountInfoFromToken function return Account from token
func GetAccountInfoFromToken(token string) (Account, error) {
	userClaims, err := GetUserClaimsFromToken(token)
	if err != nil {
		return Account{}, err
	}

	return ToAccount(userClaims), nil
}

// GetLoginInfo function return userID, companyID, branchID, departmentID
func GetLoginInfo(ctx context.Context) (int64, int64, int64, int64, error) {
	accessToken, suffix, err := GetLoginAccessToken(ctx)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	var companyID, branchID, departmentID int64
	if len(suffix) > 0 {
		parts := strings.Split(suffix, "|")
		if len(parts) >= 3 {
			companyID, _ = ToInt64(parts[0])
			branchID, _ = ToInt64(parts[1])
			departmentID, _ = ToInt64(parts[2])
		}
	}

	userClaims, err := GetUserClaimsFromToken(accessToken)
	if err != nil {
		return 0, 0, 0, 0, Unauthenticated
	}
	userID, err := ToInt64((*userClaims)["userId"])
	if err != nil {
		return 0, 0, 0, 0, Unauthenticated
	}

	return userID, companyID, branchID, departmentID, nil
}

// GetLoginInfoV2 function return struct { UserID, CompanyID, BranchID, DepartmentID}
func GetLoginInfoV2(ctx context.Context) (LoginInfo, error) {
	accessToken, suffix, err := GetLoginAccessToken(ctx)
	if err != nil {
		return LoginInfo{}, err
	}

	var companyID, branchID, departmentID int64
	var diffHour float64 = DiffHourNil
	if len(suffix) > 0 {
		parts := strings.Split(suffix, "|")
		if len(parts) >= 3 {
			companyID, _ = ToInt64(parts[0])
			branchID, _ = ToInt64(parts[1])
			departmentID, _ = ToInt64(parts[2])

			if len(parts) >= 4 {
				diffHour = ToFloat64WithDefault(parts[3], DiffHourNil)
			}
		}
	}

	userClaims, err := GetUserClaimsFromToken(accessToken)
	if err != nil {
		return LoginInfo{}, Unauthenticated
	}
	userID, err := ToInt64((*userClaims)["userId"])
	username, _ := ToString((*userClaims)["username"])
	if err != nil {
		return LoginInfo{}, Unauthenticated
	}
	partnerId, _ := ToInt64((*userClaims)["partnerId"])
	partnerCode, _ := ToString((*userClaims)["partnerCode"])
	partnerName, _ := ToString((*userClaims)["partnerName"])
	ip, _ := ToString((*userClaims)["ip"])
	deviceId, _ := ToInt64((*userClaims)["deviceId"])
	accountType, _ := ToInt32((*userClaims)["accountType"])

	return LoginInfo{
		UserId:       userID,
		Username:     username,
		CompanyId:    companyID,
		BranchId:     branchID,
		DepartmentId: departmentID,
		PartnerId:    partnerId,
		PartnerCode:  partnerCode,
		Ip:           ip,
		PartnerName:  partnerName,
		DiffHour:     diffHour,
		DeviceId:     deviceId,
		AccountType:  accountType,
	}, nil
}

func DecodeToken(accessToken string) (LoginInfo, error) {
	_, suffix, err := getTokenSuffix(accessToken)

	if err != nil {
		return LoginInfo{}, err
	}

	var companyID, branchID, departmentID int64
	var diffHour float64 = DiffHourNil
	if len(suffix) > 0 {
		parts := strings.Split(suffix, "|")
		if len(parts) >= 3 {
			companyID, _ = ToInt64(parts[0])
			branchID, _ = ToInt64(parts[1])
			departmentID, _ = ToInt64(parts[2])

			if len(parts) >= 4 {
				diffHour = ToFloat64WithDefault(parts[3], DiffHourNil)
			}
		}
	}

	userClaims, err := GetUserClaimsFromToken(accessToken)
	if err != nil {
		return LoginInfo{}, Unauthenticated
	}
	userID, err := ToInt64((*userClaims)["userId"])
	username, _ := ToString((*userClaims)["username"])
	if err != nil {
		return LoginInfo{}, Unauthenticated
	}
	partnerId, _ := ToInt64((*userClaims)["partnerId"])
	partnerCode, _ := ToString((*userClaims)["partnerCode"])
	partnerName, _ := ToString((*userClaims)["partnerName"])
	deviceId, _ := ToInt64((*userClaims)["deviceId"])
	accountType, _ := ToInt32((*userClaims)["accountType"])

	return LoginInfo{
		UserId:       userID,
		Username:     username,
		CompanyId:    companyID,
		BranchId:     branchID,
		DepartmentId: departmentID,
		PartnerId:    partnerId,
		PartnerCode:  partnerCode,
		PartnerName:  partnerName,
		DiffHour:     diffHour,
		DeviceId:     deviceId,
		AccountType:  accountType,
	}, nil
}

// GetUserID function return user id from context
func GetUserID(ctx context.Context) (int64, error) {
	userClaims, err := GetUserClaims(ctx)
	if err != nil {
		// skylog.Error(err)
		fmt.Println(err.Error())
		return 0, err
	}
	userID, _ := ToInt64((*userClaims)["userId"])
	return userID, nil
}

// GetUserIDFromToken function return user id from token
func GetUserIDFromToken(token string) (int64, error) {
	userClaims, err := GetUserClaimsFromToken(token)
	if err != nil {
		return 0, err
	}

	return ToInt64((*userClaims)["userId"])
}

// GetUserClaimsFromToken function return user claims from token
func GetUserClaimsFromToken(token string) (*jwt.MapClaims, error) {
	if len(token) == 0 {
		return nil, BadRequest
	}
	claims, err := JwtManagerInstance.Verify(token)
	if err != nil {
		// skylog.Error(err)
		fmt.Println(err.Error())
		return nil, err
	}
	if claims == nil {
		return nil, errors.New("token verification failed: no claims returned")
	}
	return claims, nil
}
