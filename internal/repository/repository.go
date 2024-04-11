package repository

import (
	"context"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
)

type ChatRepository interface {
	CreateChat(ctx context.Context) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
	AddUsersToChat(ctx context.Context, chatID int64, usernames []string) error
	CreateMessage(ctx context.Context, chatID int64, sender string, text string) (int64, error)
	GetChatMessages(ctx context.Context, chatID, page, pageSize int64) ([]*serviceModel.ChatMessage, error)
	GetChatMessagesCount(ctx context.Context, chatID int64) (uint64, error)
}
