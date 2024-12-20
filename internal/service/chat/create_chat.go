package chat

import (
	"context"
)

func (s *chatService) CreateChat(ctx context.Context, usernames []string, title string) (int64, error) {
	err := s.userClient.CheckUsersExistence(ctx, usernames)
	if err != nil {
		return 0, err
	}

	var id int64

	err = s.txManager.ReadCommitted(
		ctx, func(ctx context.Context) error {
			chatID, err := s.repo.CreateChat(ctx, title)
			if err != nil {
				return err
			}

			err = s.repo.AddUsersToChat(ctx, chatID, usernames)
			if err != nil {
				return err
			}

			id = chatID

			return nil
		},
	)

	if err != nil {
		return 0, err
	}

	s.InitMessagesChan(id)

	return id, nil
}
