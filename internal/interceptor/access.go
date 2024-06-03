package interceptor

import (
	"context"
	descAccess "github.com/gomscourse/auth/pkg/access_v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GetAccessInterceptor(client descAccess.AccessV1Client) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("metadata is not provided")
		}

		ctx = metadata.NewOutgoingContext(ctx, md)

		_, err := client.Check(
			ctx, &descAccess.CheckRequest{
				EndpointAddress: info.FullMethod,
			},
		)

		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
