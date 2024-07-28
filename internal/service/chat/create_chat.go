package chat

import (
	"context"
	"github.com/gomscourse/chat-server/internal/logger"
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

	messagesCh := make(chan *serviceModel.ChatMessage)
	s.channels[id] = messagesCh

	go func() {
		for {
			select {
			case msg, okChan := <-messagesCh:
				if !okChan {
					return
				}

				for _, st := range s.chats[id].streams {
					if err := st.Send(msg); err != nil {
						logger.Error(err.Error(), "chatID", id)
					}
				}
			}
		}
	}()

	return id, nil
}
