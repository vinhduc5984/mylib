package utils

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Account struct {
	Id             int64   `json:"id"`
	PartnerId      int64   `json:"partnerId"`
	PartnerCode    string  `json:"partnerCode"`
	PartnerName    string  `json:"partnerName"`
	DiffHour       float64 `json:"diffHour"`
	Username       *string `json:"username" orm:"null"`
	FullName       string  `json:"fullName"`
	DeviceId       int64   `json:"deviceId"`
	AccountType    int32   `json:"accountType"`
	JobTitle       string  `json:"jobTitle"`
	MemberOfDeptId int64   `json:"memberOfDeptId"`
	Ip             string  `json:"ip"`
}

// JwtManager struct
type JwtManager struct {
	VerifyKey *rsa.PublicKey
	SignKey   *rsa.PrivateKey
}

// JwtManager global instance
var (
	JwtManagerInstance *JwtManager
)

// UserClaims struct holds custom jwt claim
type UserClaims struct {
	jwt.StandardClaims
	UserID      string `json:"userId"`
	PartnerId   string `json:"partnerId"`
	PartnerName string `json:"partnerName"`
	DiffHour    string `json:"diffHour"`
	Username    string `json:"username"`
	FullName    string `json:"fullName"`
	DeviceId    string `json:"deviceId"`
	AccountType string `json:"accountType"`
	Ip          string `json:"ip"`
}

// NewJwtManager function return new *JwtManager
func NewJwtManager(signKey *rsa.PrivateKey, verifyKey *rsa.PublicKey) *JwtManager {
	return &JwtManager{
		SignKey:   signKey,
		VerifyKey: verifyKey,
	}
}

// Generate function return token of account
func (manager *JwtManager) Generate(isRefreshToken bool, account Account, expDuration time.Duration) (string, error) {
	var claims UserClaims

	if isRefreshToken {
		claims = UserClaims{
			UserID:      fmt.Sprintf("%v", account.Id),
			PartnerId:   fmt.Sprintf("%v", account.PartnerId),
			PartnerName: account.PartnerName,
			DiffHour:    fmt.Sprintf("%v", account.DiffHour),
			Username:    *account.Username,
			FullName:    account.FullName,
			DeviceId:    fmt.Sprintf("%v", account.DeviceId),
			AccountType: fmt.Sprintf("%v", account.AccountType),
			Ip:          fmt.Sprintf("%v", account.Ip),
		}
	} else {
		claims = UserClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(expDuration).Unix(),
			},
			UserID:      fmt.Sprintf("%v", account.Id),
			PartnerId:   fmt.Sprintf("%v", account.PartnerId),
			PartnerName: account.PartnerName,
			DiffHour:    fmt.Sprintf("%v", account.DiffHour),
			Username:    *account.Username,
			FullName:    account.FullName,
			DeviceId:    fmt.Sprintf("%v", account.DeviceId),
			AccountType: fmt.Sprintf("%v", account.AccountType),
			Ip:          fmt.Sprintf("%v", account.Ip),
		}
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	return token.SignedString(manager.SignKey)
}

// Verify function return UserClaims of given token
func (manager *JwtManager) Verify(token string) (*jwt.MapClaims, error) {
	index := strings.Index(token, "|||")
	if index >= 0 {
		token = token[:index]
	}
	j, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if manager == nil {
			return nil, errors.New("JwtManager is nil")
		}
		if manager.VerifyKey == nil {
			return nil, errors.New("JwtManager.VerifyKey is nil")
		}
		return manager.VerifyKey, nil
	})
	if err != nil {
		if errorValue, ok := err.(*jwt.ValidationError); ok {
			switch errorValue.Errors {
			case jwt.ValidationErrorExpired:
				return nil, errors.New("SYS.MSG.VALIDATION_EXPIRED_ERROR")
			case jwt.ValidationErrorMalformed:
				return nil, errors.New("SYS.MSG.VALIDATION_MALFORMED_ERROR")
			}
		}
		return nil, err
	}
	claims, ok := j.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid JWT claims format")
	}
	return &claims, nil
}

func (manager *JwtManager) MakeToken(data map[string]interface{}, expDuration time.Duration) (string, error) {
	claims := &jwt.MapClaims{
		"iss":  "issuer",
		"exp":  time.Now().Add(expDuration).Unix(),
		"data": data,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	return token.SignedString(manager.SignKey)
}

func (manager *JwtManager) ParseToken(tokenString string) (map[string]interface{}, error) {
	j, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return manager.VerifyKey, nil
	})

	if err != nil {
		errorValue, _ := err.(*jwt.ValidationError)
		switch errorValue.Errors {
		case jwt.ValidationErrorExpired:
			return nil, errors.New("SYS.MSG.VALIDATION_EXPIRED_ERROR")
		case jwt.ValidationErrorMalformed:
			return nil, errors.New("SYS.MSG.VALIDATION_MALFORMED_ERROR")
		}
		return nil, err
	}
	claims := j.Claims.(jwt.MapClaims)
	data := claims["data"].(map[string]interface{})
	return data, nil
}
