package app

import (
	"context"
	chatApi "github.com/gomscourse/chat-server/internal/api/chat"
	"github.com/gomscourse/chat-server/internal/client/db"
	"github.com/gomscourse/chat-server/internal/client/db/pg"
	"github.com/gomscourse/chat-server/internal/closer"
	"github.com/gomscourse/chat-server/internal/config"
	"github.com/gomscourse/chat-server/internal/config/env"
	"github.com/gomscourse/chat-server/internal/repository"
	chatRepo "github.com/gomscourse/chat-server/internal/repository/chat"
	"github.com/gomscourse/chat-server/internal/service"
	chatService "github.com/gomscourse/chat-server/internal/service/chat"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	pgPool         *pgxpool.Pool
	dbClient       db.Client
	chatRepository repository.ChatRepository
	chatService    service.ChatService
	chatImpl       *chatApi.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) PgConfig() config.PGConfig {
	if sp.pgConfig == nil {
		pgConfig, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to initialize PG config: %s", err.Error())
		}

		sp.pgConfig = pgConfig
	}

	return sp.pgConfig
}

func (sp *serviceProvider) GRPCConfig() config.GRPCConfig {
	if sp.grpcConfig == nil {
		grpcConfig, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to initialize GRPC config: %s", err.Error())
		}

		sp.grpcConfig = grpcConfig
	}

	return sp.grpcConfig
}

func (sp *serviceProvider) PGPool(ctx context.Context) *pgxpool.Pool {
	if sp.pgPool == nil {
		pool, err := pgxpool.Connect(ctx, sp.PgConfig().DSN())
		if err != nil {
			log.Fatalf("failed to initialize PG pool: %s", err.Error())
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		sp.pgPool = pool
	}

	return sp.pgPool
}

func (sp *serviceProvider) DBClient(ctx context.Context) db.Client {
	if sp.dbClient == nil {
		client, err := pg.New(ctx, sp.PgConfig().DSN())
		if err != nil {
			log.Fatalf("failed to initialize DB client: %s", err.Error())
		}

		sp.dbClient = client
	}

	return sp.dbClient
}

func (sp *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if sp.chatRepository == nil {
		sp.chatRepository = chatRepo.NewChatRepository(sp.DBClient(ctx))
	}

	return sp.chatRepository
}

func (sp *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if sp.chatService == nil {
		sp.chatService = chatService.NewChatService(sp.ChatRepository(ctx))
	}

	return sp.chatService
}

func (sp *serviceProvider) ChatImpl(ctx context.Context) *chatApi.Implementation {
	if sp.chatImpl == nil {
		sp.chatImpl = chatApi.NewImplementation(sp.ChatService(ctx))
	}

	return sp.chatImpl
}
