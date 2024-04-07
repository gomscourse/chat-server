package repository

import "context"

type ChatRepository interface {
	CreateChat(ctx context.Context) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
	AddUsersToChat(ctx context.Context, chatID int64, usernames []string) error
	CreateMessage(ctx context.Context, chatID int64, sender string, text string) (int64, error)
}
