package repository

import (
	"context"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, title string) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
	AddUsersToChat(ctx context.Context, chatID int64, usernames []string) error
	CreateMessage(ctx context.Context, chatID int64, sender string, text string) (*serviceModel.ChatMessage, error)
	GetChatMessages(ctx context.Context, chatID, page, pageSize int64) ([]*serviceModel.ChatMessage, error)
	GetChatMessagesCount(ctx context.Context, chatID int64) (uint64, error)
	CheckUserChat(ctx context.Context, chatID int64, username string) (bool, error)
}
