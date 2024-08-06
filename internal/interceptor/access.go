package interceptor

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	descAccess "github.com/gomscourse/auth/pkg/access_v1"
	"github.com/gomscourse/chat-server/internal/context_keys"
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

const authPrefix = "Bearer "

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     int32  `json:"role"`
}

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
		span, ctx := opentracing.StartSpanFromContext(ctx, "check access")
		defer span.Finish()
		err := checkAccess(ctx, client, info.FullMethod)
		if err != nil {
			return nil, err
		}

		ctx, err = putUsernameFromJWTToContext(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func GetAccessStreamInterceptor(client descAccess.AccessV1Client) func(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := checkAccess(ss.Context(), client, info.FullMethod)
		if err != nil {
			return err
		}

		ctx, err := putUsernameFromJWTToContext(ss.Context())
		if err != nil {
			return err
		}

		wrapped := &wrappedStream{
			ctx:          ctx,
			serverStream: ss,
		}

		return handler(srv, wrapped)
	}
}

func checkAccess(ctx context.Context, client descAccess.AccessV1Client, method string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return sys.NewCommonError("metadata is not provided", codes.InvalidArgument)
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := client.Check(
		ctx, &descAccess.CheckRequest{
			EndpointAddress: method,
		},
	)

	return err
}

func putUsernameFromJWTToContext(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, &UserClaims{})
	if err != nil {
		return nil, sys.NewCommonError("failed to parse token claims", codes.InvalidArgument)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, sys.NewCommonError("invalid user claims", codes.InvalidArgument)
	}

	return context.WithValue(ctx, context_keys.UsernameKey, claims.Username), nil
}

type wrappedStream struct {
	serverStream grpc.ServerStream
	ctx          context.Context
}

func (ws *wrappedStream) SetHeader(md metadata.MD) error {
	return ws.serverStream.SetHeader(md)
}
func (ws *wrappedStream) SendHeader(md metadata.MD) error {
	return ws.serverStream.SendHeader(md)
}
func (ws *wrappedStream) SetTrailer(md metadata.MD) {
	ws.serverStream.SetTrailer(md)
}
func (ws *wrappedStream) Context() context.Context {
	return ws.ctx
}
func (ws *wrappedStream) SendMsg(m any) error {
	return ws.serverStream.SendMsg(m)
}
func (ws *wrappedStream) RecvMsg(m any) error {
	return ws.serverStream.RecvMsg(m)
}
