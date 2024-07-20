package service

import (
	"context"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
)

type ChatService interface {
	CreateChat(ctx context.Context, usernames []string, title string) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) error
	SendMessage(ctx context.Context, sender, text string, chatID int64) error
	GetChatMessages(ctx context.Context, chatID, page, pageSize int64) ([]*serviceModel.ChatMessage, error)
	GetChatMessagesCount(ctx context.Context, chatID int64) (uint64, error)
}
