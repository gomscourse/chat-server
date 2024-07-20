package chat

import (
	"context"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
)

func (s *chatService) CreateChat(ctx context.Context, usernames []string, title string) (int64, error) {
	var id int64

	err := s.txManager.ReadCommitted(
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

	s.channels[id] = make(chan *serviceModel.ChatMessage)

	//TODO: прослушивание сообщений из канала чата и отправка их через стримы

	return id, nil
}
