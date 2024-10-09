package env

import (
	"github.com/gomscourse/chat-server/internal/config"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	grpcHostEnvName             = "GRPC_HOST"
	grpcPortEnvName             = "GRPC_PORT"
	grpcAccessClientHostEnvName = "GRPC_HOST_ACCESS_CLIENT"
	grpcAccessClientPortEnvName = "GRPC_PORT_ACCESS_CLIENT"
)

type grpcConfig struct {
	host              string
	port              string
	accessClientHost  string
	accessClientPort  string
	requestLimitCount int
	requestLimitTime  time.Duration
}

func NewGRPCConfig() (config.GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	accessClientHost := os.Getenv(grpcAccessClientHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("access client grpc host not found")
	}

	accessClientPort := os.Getenv(grpcAccessClientPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("access client grpc port not found")
	}

	return &grpcConfig{
		host:              host,
		port:              port,
		accessClientHost:  accessClientHost,
		accessClientPort:  accessClientPort,
		requestLimitCount: 100,
		requestLimitTime:  time.Second,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *grpcConfig) AuthServiceAddress() string {
	return net.JoinHostPort(cfg.accessClientHost, cfg.accessClientPort)
}

func (cfg *grpcConfig) RateLimit() (int, time.Duration) {
	return cfg.requestLimitCount, cfg.requestLimitTime
}
