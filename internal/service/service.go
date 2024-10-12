package service

import (
	"context"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
)

type UserClient interface {
	CheckUsersExistence(ctx context.Context, usernames []string) error
}

type ChatService interface {
	CreateChat(ctx context.Context, usernames []string, title string) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) error
	SendMessage(ctx context.Context, text string, chatID int64) error
	GetChatMessagesAndCount(ctx context.Context, chatID, page, pageSize int64) (
		[]*serviceModel.ChatMessage,
		uint64,
		error,
	)
	GetAvailableChatsAndCount(
		ctx context.Context, page, pageSize int64,
	) ([]*serviceModel.Chat, uint64, error)
	ConnectChat(stream Stream, chatID int64) error
	InitMessagesChan(chatID int64) chan *serviceModel.ChatMessage
	GetChannels() map[int64]chan *serviceModel.ChatMessage
	CheckChatAvailability(ctx context.Context, chatID int64, username string) error
	CheckCtxUserChatAvailability(ctx context.Context, chatID int64) error
}

type Stream interface {
	Send(message *serviceModel.ChatMessage) error
	Context() context.Context
}
