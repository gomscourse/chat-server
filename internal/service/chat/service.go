package chat

import (
	"github.com/gomscourse/chat-server/internal/model"
	"github.com/gomscourse/chat-server/internal/repository"
	"github.com/gomscourse/chat-server/internal/service"
	"github.com/gomscourse/common/pkg/db"
	"sync"
)

type Chat struct {
	streams map[string]service.Stream
	m       sync.RWMutex
}

type chatService struct {
	repo       repository.ChatRepository
	txManager  db.TxManager
	userClient service.UserClient

	channels  map[int64]chan *model.ChatMessage
	mxChannel sync.RWMutex

	chats  map[int64]*Chat
	mxChat sync.RWMutex
}

func NewChatService(
	repo repository.ChatRepository,
	manager db.TxManager,
	userClient service.UserClient,
) service.ChatService {
	return &chatService{
		repo:       repo,
		txManager:  manager,
		userClient: userClient,
		channels:   make(map[int64]chan *model.ChatMessage),
		chats:      make(map[int64]*Chat),
	}
}

func NewTestService(deps ...interface{}) service.ChatService {
	srv := chatService{
		channels: make(map[int64]chan *model.ChatMessage),
		chats:    make(map[int64]*Chat),
	}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.ChatRepository:
			srv.repo = s
		case db.TxManager:
			srv.txManager = s
		}
	}

	return &srv
}
