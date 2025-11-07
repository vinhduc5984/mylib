package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/avct/uasurfer"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
	// "suntech.com.vn/skylib/skylog.git/skylog"
)

var (
	GlobalAuthInterceptor *AuthInterceptor
)

// AuthInterceptor struct
type AuthInterceptor struct {
	jwtManager *JwtManager
}

// NewAuthInterceptor function: create new AuthInterceptor
func NewAuthInterceptor(jwtManager *JwtManager) *AuthInterceptor {

	return &AuthInterceptor{
		jwtManager: jwtManager,
	}
}

// Init function
func Init(jwtManager *JwtManager) {
	GlobalAuthInterceptor = NewAuthInterceptor(jwtManager)
}

// Unary interceptor function
func (interceptor *AuthInterceptor) Unary(publicMethods map[string]bool, lockScreens map[int64]bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// skylog.Info("Unary Func publicMethods", publicMethods)
		if err := interceptor.authorize(ctx, info.FullMethod, publicMethods, lockScreens); err != nil {
			// skylog.Error(err)
			return nil, err
		}
		// skylog.Info("Unary Func publicMethods authorize done")
		PrintRequest(info.FullMethod, req)
		return handler(ctx, req)
	}
}

// ClientUnary interceptor function
func (interceptor *AuthInterceptor) ClientUnary() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
	}
}

// ClientUnaryWithPackage interceptor function
func (interceptor *AuthInterceptor) ClientUnaryWithPackage(test, replaceWith string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		method = strings.ReplaceAll(method, test, replaceWith)
		return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
	}
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	auths := md["authorization"]
	if len(auths) == 0 {
		return nil
	}
	accessToken := strings.Replace(auths[0], "Bearer ", "", 1)
	params := []string{
		"authorization",
		accessToken,
	}
	if grpcGatewayOrigin, ok := md["grpcgateway-origin"]; ok {
		params = append(params, "grpcgateway-origin", grpcGatewayOrigin[0])
	}
	if origin, ok := md["origin"]; ok {
		params = append(params, "origin", origin[0])
	}
	if xForwardedHost, ok := md["x-forwarded-host"]; ok {
		params = append(params, "x-forwarded-host", xForwardedHost[0])
	}

	return metadata.AppendToOutgoingContext(ctx, params...)
}

// Stream interceptor function
func (interceptor *AuthInterceptor) Stream(publicMethods map[string]bool, lockScreens map[int64]bool) grpc.StreamServerInterceptor {
	return func(server interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// skylog.Info("Stream function publicMethods 1", publicMethods)
		if err := interceptor.authorize(stream.Context(), info.FullMethod, publicMethods, lockScreens); err != nil {
			// skylog.Error(err)
			return err
		}
		// skylog.Info("Stream function publicMethods 2")
		PrintRequest(info.FullMethod, "")
		return handler(server, stream)
	}
}

// ClientStream interceptor function
func (interceptor *AuthInterceptor) ClientStream() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(interceptor.attachToken(ctx), desc, cc, method, opts...)
	}
}

// ClientStream interceptor function
func (interceptor *AuthInterceptor) ClientStreamWithPackage(test, replaceWith string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		method = strings.ReplaceAll(method, test, replaceWith)
		return streamer(interceptor.attachToken(ctx), desc, cc, method, opts...)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string, publicMethods map[string]bool, lockScreens map[int64]bool) error {
	// skylog.Info(method)
	// skylog.Info(publicMethods[method])
	if publicMethods[method] {
		return nil
	}
	// skylog.Info("userID 1")
	userID, err := GetUserID(ctx)
	if err != nil {
		// skylog.Error(err)
		return err
	}
	// skylog.Info("userID 2", userID)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ua := uasurfer.Parse(md["user-agent"][0])
		fmt.Println(lockScreens[userID])
		if ua.DeviceType == uasurfer.DeviceComputer && lockScreens[userID] {
			return NeedLogin
		}
	}

	return err
}
