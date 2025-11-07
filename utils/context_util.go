package utils

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func MakeContext(r *http.Request) context.Context {
	// Extract headers from the HTTP request
	md := metadata.New(nil)
	md.Append("pattern", r.URL.String())
	md.Append("x-forwarded-host", r.Host)

	for key, values := range r.Header {
		for _, value := range values {
			md.Append(key, value)
		}
	}

	// Create a gRPC context with the metadata
	ctx := metadata.NewIncomingContext(context.Background(), md)

	// Optionally, add other values to the context if needed
	ctx = context.WithValue(ctx, "http-method", r.Method)

	return ctx
}
