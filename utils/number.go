package utils

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

var (
	base58Chars = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
)

// ToInt64 convert interface{} to int64.
func ToInt64(source interface{}) (int64, error) {
	if source == nil {
		return 0, nil
	}

	str := fmt.Sprintf("%v", source)

	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return -1, err
	}
	return num, nil
}

// ToFloat64 convert interface{} to float64.
func ToFloat64(source interface{}) (float64, error) {
	if source == nil {
		return 0, nil
	}

	str := fmt.Sprintf("%v", source)

	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func ToFloat64WithDefault(source interface{}, defaultVal float64) float64 {
	if source == nil {
		return defaultVal
	}

	str := fmt.Sprintf("%v", source)

	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return defaultVal
	}
	return num
}

// ToF64 convert interface{} to float64.
func ToF64(source interface{}) float64 {
	num, err := ToFloat64(source)
	if err != nil {
		return 0.0
	}

	return num
}

// ToF64Ptr convert interface{} to *float64.
func ToF64Ptr(source interface{}) *float64 {
	num, err := ToFloat64(source)
	if err != nil {
		return nil
	}

	return &num
}

// ToI64 convert interface{} to int64.
func ToI64(source interface{}) int64 {
	num, err := ToInt64(source)
	if err != nil {
		return 0
	}

	return num
}

// ToI64Ptr convert interface{} to *int64.
func ToI64Ptr(source interface{}) *int64 {
	num, err := ToInt64(source)
	if err != nil {
		return nil
	}

	return &num
}

// ToInt32 convert interface{} to int32.
func ToInt32(source interface{}) (int32, error) {
	num, err := ToInt64(source)
	if err != nil {
		return 0, err
	}

	return int32(num), nil
}

// ToI32 convert interface{} to int32.
func ToI32(source interface{}) int32 {
	num, err := ToInt32(source)
	if err != nil {
		return 0
	}

	return num
}

// ToI32Ptr convert interface{} to *int32.
func ToI32Ptr(source interface{}) *int32 {
	num, err := ToInt32(source)
	if err != nil {
		return nil
	}

	return &num
}

// ToInt16 convert interface{} to int16.
func ToInt16(source interface{}) (int16, error) {
	num, err := ToInt64(source)
	if err != nil {
		return 0, err
	}

	return int16(num), nil
}

// ToI16 convert interface{} to int16.
func ToI16(source interface{}) int16 {
	num, err := ToInt16(source)
	if err != nil {
		return 0
	}

	return num
}

// ToI16Ptr convert interface{} to *int16.
func ToI16Ptr(source interface{}) *int16 {
	num, err := ToInt16(source)
	if err != nil {
		return nil
	}

	return &num
}

func Round(val float64, precision int) float64 {
	return math.Round(val*(math.Pow10(precision))) / math.Pow10(precision)
}

func ZipNumber(val int64) int {
	var sum int64
	countLoop := 0
	for val >= 10 {
		sum = 0

		for _, char := range strings.Split(strconv.FormatInt(val, 10), "") {
			sum += ToI64(char)
		}

		countLoop++
		val = sum
	}
	return int(val)
}

func Base10ToBase58(val *big.Int) string {
	var buffer bytes.Buffer
	var remainder *big.Int

	zero := big.NewInt(0)
	int58 := big.NewInt(58)

	for val.Cmp(zero) > 0 {
		val, remainder = new(big.Int).DivMod(val, int58, new(big.Int))
		buffer.WriteByte(base58Chars[remainder.Int64()])
	}

	// Reverse the buffer to get the correct order
	result := buffer.Bytes()
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

func Base58ToBase10(val string) *big.Int {
	base58CharsStr := string(base58Chars)

	var tmpVal *big.Int
	num := big.NewInt(0)
	for i, char := range strings.Split(val, "") {
		pos := strings.IndexByte(base58CharsStr, char[0])
		if pos == -1 {
			return big.NewInt(0)
		}
		tmpVal = big.NewInt(58)
		tmpVal.Exp(tmpVal, big.NewInt(int64(len(val)-i-1)), nil)
		tmpVal.Mul(tmpVal, big.NewInt(int64(pos)))
		num.Add(num, tmpVal)
	}
	return num
}

func ToInt(val interface{}, defaultValue int) int {
	switch v := val.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		num, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return int(num)
		} else {
			return defaultValue
		}
	default:
		num, err := strconv.ParseInt(ToStr(val), 10, 64)
		if err == nil {
			return int(num)
		}
		return defaultValue
	}
}
