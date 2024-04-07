package chat

import (
	"context"
)

func (s chatService) CreateChat(ctx context.Context, usernames []string) (int64, error) {
	//TODO: обернуть в транзакцию
	chatID, err := s.repo.CreateChat(ctx)
	if err != nil {
		return 0, err
	}

	err = s.repo.AddUsersToChat(ctx, chatID, usernames)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
