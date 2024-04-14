package chat

import (
	"github.com/gomscourse/chat-server/internal/repository"
	"github.com/gomscourse/chat-server/internal/service"
	"github.com/gomscourse/common/pkg/db"
)

type chatService struct {
	repo      repository.ChatRepository
	txManager db.TxManager
}

func NewChatService(repo repository.ChatRepository, manager db.TxManager) service.ChatService {
	return &chatService{repo: repo, txManager: manager}
}
