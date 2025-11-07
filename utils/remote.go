package utils

import (
	"context"

	"google.golang.org/grpc"
)

func MakeRemoteConn(remoteHost string) (*grpc.ClientConn, error) {
	return grpc.Dial(remoteHost, grpc.WithInsecure(), grpc.WithUnaryInterceptor(GlobalAuthInterceptor.ClientUnary()), grpc.WithStreamInterceptor(GlobalAuthInterceptor.ClientStream()))
}

func MakeRemoteConnWithPackage(remoteHost, test, replaceWith string) (*grpc.ClientConn, error) {
	return grpc.Dial(remoteHost, grpc.WithInsecure(), grpc.WithUnaryInterceptor(GlobalAuthInterceptor.ClientUnaryWithPackage(test, replaceWith)), grpc.WithStreamInterceptor(GlobalAuthInterceptor.ClientStreamWithPackage(test, replaceWith)))
}

func MakeRemoteConnContextWithPackage(ctx context.Context, remoteHost, test, replaceWith string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, remoteHost, grpc.WithInsecure(), grpc.WithUnaryInterceptor(GlobalAuthInterceptor.ClientUnaryWithPackage(test, replaceWith)), grpc.WithStreamInterceptor(GlobalAuthInterceptor.ClientStreamWithPackage(test, replaceWith)))
}
