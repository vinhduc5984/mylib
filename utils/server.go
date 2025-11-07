package utils

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

// Accept
// Accept-Encoding
// Accept-Language
// Cache-Control
// Connection
// Priority
// Sec-Ch-Ua
// Sec-Ch-Ua-Mobile
// Sec-Ch-Ua-Platform
// Sec-Fetch-Dest
// Sec-Fetch-Mode
// Sec-Fetch-Site
// Upgrade-Insecure-Requests
// User-Agent
// X-Forwarded-For
// X-Forwarded-Port
// X-Forwarded-Scheme
// X-Forwarded-Host
// X-Real-Ip
var defaultAcceptKey = []string{"X-Forwarded-Scheme", "X-Forwarded-Host", "X-Forwarded-For"}

func buildAcceptHeaderKeys(otherHeaderKeys ...string) map[string]bool {
	acceptHeaderKeys := make(map[string]bool)
	for _, v := range defaultAcceptKey {
		acceptHeaderKeys[strings.ToLower(v)] = true
	}
	for _, v := range otherHeaderKeys {
		_, ok := acceptHeaderKeys[strings.ToLower(v)]
		if !ok {
			acceptHeaderKeys[strings.ToLower(v)] = true
		}
	}
	return acceptHeaderKeys
}

// / create new server mux
func NewServeMux(otherHeaderKeys ...string) *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			// build all accept header key
			acceptHeaderKeys := buildAcceptHeaderKeys(otherHeaderKeys...)

			// build metadata
			md := make(map[string]string)

			// pass request header key to grpc metadata
			var headerKey string
			var ok bool
			for key, values := range r.Header {
				headerKey = strings.ToLower(key)
				_, ok = acceptHeaderKeys[headerKey]
				if ok && len(values) > 0 {
					md[headerKey] = values[0]
				}
			}

			// add request url to grpc
			md["pattern"] = r.URL.String()

			return metadata.New(md)
		}),
	)
}
