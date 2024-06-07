package app

import (
	"context"
	"flag"
	"github.com/gomscourse/chat-server/internal/config"
	"github.com/gomscourse/chat-server/internal/interceptor"
	"github.com/gomscourse/chat-server/internal/logger"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
	"github.com/gomscourse/common/pkg/closer"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/natefinch/lumberjack"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
)

var logLevel = flag.String("l", "info", "log level")

type initializer func(ctx context.Context) error

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	logger.Init(getLogHandler())
	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []initializer{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.LogInterceptor,
				interceptor.GetAccessInterceptor(a.serviceProvider.AccessClient()),
			),
		),
	)
	reflection.Register(a.grpcServer)
	desc.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatImpl(ctx))
	return nil
}

func (a *App) runGRPCServer() error {
	address := a.serviceProvider.GRPCConfig().Address()
	log.Printf("GRPC server listening at %v", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	if err = a.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func getLogLevel() slog.Level {
	var level slog.Level
	if err := level.UnmarshalText([]byte(*logLevel)); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return level
}

func getLogHandler() slog.Handler {
	return slog.NewJSONHandler(
		io.MultiWriter(
			os.Stdout,
			&lumberjack.Logger{
				Filename:   "logs/app.log",
				MaxSize:    10, // megabytes
				MaxBackups: 3,
				MaxAge:     7, // days
			},
		), &slog.HandlerOptions{
			Level: getLogLevel(),
		},
	)
}
