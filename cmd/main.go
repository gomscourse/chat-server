package main

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gomscourse/chat-server/internal/config"
	"github.com/gomscourse/chat-server/internal/config/env"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type server struct {
	desc.UnimplementedChatV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	fmt.Printf("%+v\n", req.GetUsernames())
	return &desc.CreateResponse{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	deleteBuilder := sq.Delete("chat").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": req.GetId()})
	query, args, err := deleteBuilder.ToSql()

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to delete chat: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	fmt.Printf("%+v\n", req.GetText())
	return &emptypb.Empty{}, nil
}

func main() {
	ctx := context.Background()
	// Считываем переменные окружения
	err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
