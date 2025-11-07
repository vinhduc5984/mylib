package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/mileusna/useragent"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"google.golang.org/grpc/metadata"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func UnAccent(str string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, str)
	result = strings.Replace(result, "Đ", "D", -1)
	result = strings.Replace(result, "Ð", "D", -1)
	result = strings.Replace(result, "đ", "d", -1)

	result = strings.Replace(result, "Ư", "U", -1)
	result = strings.Replace(result, "ư", "u", -1)
	result = strings.Replace(result, "Ứ", "U", -1)
	result = strings.Replace(result, "ứ", "u", -1)
	result = strings.Replace(result, "Ừ", "U", -1)
	result = strings.Replace(result, "ừ", "u", -1)
	result = strings.Replace(result, "Ự", "U", -1)
	result = strings.Replace(result, "ự", "u", -1)
	return result
}

func LowerUnAccent(str string) string {
	return strings.ToLower(UnAccent(str))
}

func UpperUnAccent(str string) string {
	return strings.ToUpper(UnAccent(str))
}

func StringPadding(input string, padLength int, padString string, padType string) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] + input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
	}

	return output
}

func StringRightPaddingList(input []string, padLength []int) string {
	result := []string{}
	for i, str := range input {
		result = append(result, StringPadding(str, padLength[i], " ", "RIGHT"))
	}

	return strings.Join(result, "")
}

func JsonPrettyAny(v interface{}) string {
	in, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("JsonPrettyAny error: %v", err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return string(in)
	}
	return out.String()
}

func WildCardLike(query string) string {
	return WildCardLikeSensitive(query, false)
}

func WildCardLikeSensitive(query string, sensitive bool) string {
	if query == "" {
		return "%"
	}
	if !sensitive {
		query = strings.ToUpper(query)
	}
	query = strings.Replace(query, "\\", "\\\\", -1)
	query = strings.Replace(query, "%", "\\%", -1)
	query = strings.Replace(query, "_", "\\_", -1)

	return "%" + query + "%"
}

func WildCardFull(query string) string {
	return WildCardFullSensitive(query, false)
}

func WildCardFullSensitive(query string, sensitive bool) string {
	if query == "" {
		return "%"
	}
	if !sensitive {
		query = strings.ToUpper(query)
	}
	query = strings.Replace(query, "\\", "\\\\", -1)
	query = strings.Replace(query, "%", "\\%", -1)
	query = strings.Replace(query, "_", "\\_", -1)

	query = strings.ReplaceAll(query, "\\s+", " ")
	query = strings.TrimSpace(query)

	arr := strings.Split(query, " ")
	query = "%"
	for _, q := range arr {
		query += q + "%"
	}

	return query
}

func GetRemoteIPAddress(r *http.Request) string {
	returnVal := ""
	if r != nil {
		// Get the client's IP address from the Request
		returnVal = r.Header.Get("X-Real-IP") // Check if the X-Real-IP header is set
		if returnVal == "" {
			returnVal = r.Header.Get("X-Forwarded-For") // Check if the X-Forwarded-For header is set
		}
		if returnVal == "" {
			returnVal = r.RemoteAddr // Fallback to RemoteAddr if headers are not set
		}

		// Extract the first IP address if X-Forwarded-For contains a comma-separated list
		if comma := strings.Index(returnVal, ","); comma >= 0 {
			returnVal = returnVal[:comma]
		}
	}
	return returnVal
}

func getOriginFromMd(md metadata.MD) string {
	if origin, success := md["grpcgateway-origin"]; success {
		if len(origin) > 0 && len(origin[0]) > 0 {
			return origin[0]
		}
	}
	if origin, success := md["origin"]; success {
		if len(origin) > 0 && len(origin[0]) > 0 {
			return origin[0]
		}
	}
	return ""
}

func GetRemoteDomainFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	// Get the client's grpcgateway from the Context
	hostName := getOriginFromMd(md)
	if len(hostName) > 0 {
		urlInfo, err := url.Parse(hostName)
		if err != nil {
			if strings.HasPrefix(hostName, "https://") {
				hostName = strings.Split(hostName, "https://")[1]
			}
			if strings.HasPrefix(hostName, "http://") {
				hostName = strings.Split(hostName, "http://")[1]
			}
			if strings.Contains(hostName, ":") {
				hostName = strings.Split(hostName, ":")[0]
			}
		} else {
			hostName = urlInfo.Hostname()
		}
		if strings.HasPrefix(hostName, "www.") {
			hostName = strings.Split(hostName, "www.")[1]
		}
		return hostName
	}

	// Get the client's host from the Context
	if forwardedHost, success := md["x-forwarded-host"]; success {
		if len(forwardedHost) > 0 && len(forwardedHost[0]) > 0 {
			hostName := forwardedHost[0]
			if strings.Contains(hostName, ":") {
				hostName = strings.Split(hostName, ":")[0]
			}
			return hostName
		}
	}

	return ""
}

func GetRemoteIPAddressFromContext(ctx context.Context) string {
	returnVal := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// Get the client's IP address from the Request
		if ip, success := md["x-real-ip"]; success {
			if len(ip) > 0 && len(ip[0]) > 0 {
				returnVal = ip[0]
			}
		}
		if returnVal == "" {
			if ip, success := md["x-forwarded-for"]; success {
				if len(ip) > 0 && len(ip[0]) > 0 {
					returnVal = ip[0]
				}
			}
		}

		// Extract the first IP address if X-Forwarded-For contains a comma-separated list
		if comma := strings.Index(returnVal, ","); comma >= 0 {
			returnVal = returnVal[:comma]
		}
	}
	return returnVal
}

func GetUserAgentInfo(ctx context.Context) (useragent.UserAgent, bool) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userAgent := ""
		if ua, hasData := md["grpcgateway-user-agent"]; hasData {
			if len(ua) > 0 && len(ua[0]) > 0 {
				userAgent = ua[0]
			}
		}
		if len(userAgent) == 0 {
			if ua, hasData := md["user-agent"]; hasData {
				if len(ua) > 0 && len(ua[0]) > 0 {
					userAgent = ua[0]
				}
			}
		}
		if len(userAgent) > 0 {
			return useragent.Parse(userAgent), true
		}
	}
	return useragent.UserAgent{}, false
}
