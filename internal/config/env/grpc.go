package env

import (
	"github.com/gomscourse/chat-server/internal/config"
	"net"
	"os"

	"github.com/pkg/errors"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type grpcConfig struct {
	host             string
	port             string
	accessClientHost string
	accessClientPort string
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

	return &grpcConfig{
		host:             host,
		port:             port,
		accessClientHost: "0.0.0.0", //TODO добавить инициализацию из .env
		accessClientPort: "50051",   //TODO добавить инициализацию из .env
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *grpcConfig) AccessClientAddress() string {
	return net.JoinHostPort(cfg.accessClientHost, cfg.accessClientPort)
}
