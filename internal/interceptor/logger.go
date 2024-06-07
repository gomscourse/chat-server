package interceptor

import (
	"context"
	"github.com/gomscourse/chat-server/internal/logger"
	"time"

	"google.golang.org/grpc"
)

func LogInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	now := time.Now()

	res, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error(), "method", info.FullMethod, "req", req)
	}

	logger.Info("request", "method", info.FullMethod, "req", req, "res", res, "duration", time.Since(now).String())

	return res, err
}
