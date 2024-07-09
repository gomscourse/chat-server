package interceptor

import (
	"context"
	"github.com/gomscourse/common/pkg/rate_limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RateLimiterInterceptor struct {
	rateLimiter *rate_limiter.TokenBucketLimiter
}

func NewRateLimiterInterceptor(limiter *rate_limiter.TokenBucketLimiter) *RateLimiterInterceptor {
	return &RateLimiterInterceptor{
		rateLimiter: limiter,
	}
}

func (rli *RateLimiterInterceptor) Unary(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if !rli.rateLimiter.Allow() {
		return nil, status.Error(codes.ResourceExhausted, "too many requests")
	}

	return handler(ctx, req)
}
