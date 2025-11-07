package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// SnakeCaseToCamelCase function
func SnakeCaseToCamelCase(inputUnderScoreStr string) (camelCase string) {
	isToUpper := false
	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}

func InsensitiveReplaceAll(source, oldValue, newValue string) string {
	re := regexp.MustCompile(`(?i)` + oldValue)
	return re.ReplaceAllString(source, newValue)
}

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	output = strings.Replace(output, "đ", "d", -1)
	output = strings.Replace(output, "Đ", "D", -1)
	return output
}

func Int64ToString(number int64) string {
	return strconv.Itoa(int(number))
}

func Int32ToString(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

// ToString convert interface{} to string.
func ToString(source interface{}) (string, error) {
	if source == nil {
		return "", nil
	}

	str := fmt.Sprintf("%v", source)

	return str, nil
}

func ToStr(val interface{}, defaultVal ...string) string {
	if val == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return ""
	}
	return fmt.Sprintf("%v", val)
}

func ReverseString(s string) string {
	runes := []rune(s) // Convert string to rune slice to handle Unicode characters properly
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i] // Swap characters
	}
	return string(runes) // Convert back to string
}

func IsEmpty(s string) bool { return len(strings.TrimSpace(s)) == 0 }

func IsNotEmpty(s string) bool { return !IsEmpty(s) }
