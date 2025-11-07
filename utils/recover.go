package utils

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	customFunc grpc_recovery.RecoveryHandlerFunc
)

type Intercepter struct {
	Unary  grpc.UnaryServerInterceptor
	Stream grpc.StreamServerInterceptor
}

// Recover panic error grpc interceptor
func RecoverInterceptor(message string) Intercepter {
	customFunc = func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "%s, panic triggered: %v", message, p)
	}
	// Shared options for the logger, with a custom gRPC code to log level function.
	rcoveryOpt := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}
	recoverUnary := grpc_recovery.UnaryServerInterceptor(rcoveryOpt...)
	recoverStream := grpc_recovery.StreamServerInterceptor(rcoveryOpt...)

	return Intercepter{
		recoverUnary,
		recoverStream,
	}
}
