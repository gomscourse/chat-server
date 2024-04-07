package chat

import (
	"github.com/gomscourse/chat-server/internal/repository"
	"github.com/gomscourse/chat-server/internal/service"
)

type chatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) service.ChatService {
	return chatService{repo: repo}
}
