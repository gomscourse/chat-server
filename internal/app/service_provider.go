package app

import (
	"context"
	descAccess "github.com/gomscourse/auth/pkg/access_v1"
	chatApi "github.com/gomscourse/chat-server/internal/api/chat"
	"github.com/gomscourse/chat-server/internal/config"
	"github.com/gomscourse/chat-server/internal/config/env"
	"github.com/gomscourse/chat-server/internal/repository"
	chatRepo "github.com/gomscourse/chat-server/internal/repository/chat"
	"github.com/gomscourse/chat-server/internal/service"
	chatService "github.com/gomscourse/chat-server/internal/service/chat"
	"github.com/gomscourse/common/pkg/closer"
	"github.com/gomscourse/common/pkg/db"
	"github.com/gomscourse/common/pkg/db/pg"
	"github.com/gomscourse/common/pkg/db/transaction"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	chatRepository repository.ChatRepository
	chatService    service.ChatService
	chatImpl       *chatApi.Implementation
	accessClient   descAccess.AccessV1Client
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

func (sp *serviceProvider) DBClient(ctx context.Context) db.Client {
	if sp.dbClient == nil {
		client, err := pg.New(ctx, sp.PgConfig().DSN())
		if err != nil {
			log.Fatalf("failed to initialize DB client: %s", err.Error())
		}

		if err = client.DB().Ping(ctx); err != nil {
			log.Fatalf("failed to ping DB: %s", err.Error())
		}

		closer.Add(client.Close)

		sp.dbClient = client
	}

	return sp.dbClient
}

func (sp *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if sp.txManager == nil {
		sp.txManager = transaction.NewTransactionManager(sp.DBClient(ctx).DB())
	}

	return sp.txManager
}

func (sp *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if sp.chatRepository == nil {
		sp.chatRepository = chatRepo.NewChatRepository(sp.DBClient(ctx))
	}

	return sp.chatRepository
}

func (sp *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if sp.chatService == nil {
		sp.chatService = chatService.NewChatService(sp.ChatRepository(ctx), sp.TxManager(ctx))
	}

	return sp.chatService
}

func (sp *serviceProvider) ChatImpl(ctx context.Context) *chatApi.Implementation {
	if sp.chatImpl == nil {
		sp.chatImpl = chatApi.NewImplementation(sp.ChatService(ctx))
	}

	return sp.chatImpl
}

func (sp *serviceProvider) AccessClient() descAccess.AccessV1Client {
	if sp.accessClient == nil {
		creds, err := credentials.NewClientTLSFromFile("service.pem", "")
		if err != nil {
			log.Fatalf("could not process the credentials: %v", err)
		}

		conn, err := grpc.Dial(
			sp.GRPCConfig().AccessClientAddress(),
			grpc.WithTransportCredentials(creds),
		)

		if err != nil {
			log.Fatalf("failed to initialize access client: %s", err.Error())
		}

		sp.accessClient = descAccess.NewAccessV1Client(conn)
	}

	return sp.accessClient
}
