package service

import "context"

type ChatService interface {
	CreateChat(ctx context.Context, usernames []string) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) error
	SendMessage(ctx context.Context, sender, text string, chatID int64) error
}
