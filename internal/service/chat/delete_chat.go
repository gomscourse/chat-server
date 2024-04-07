package chat

import "context"

func (s chatService) DeleteChat(ctx context.Context, chatID int64) error {
	return s.repo.DeleteChat(ctx, chatID)
}
