package chat

import (
	"context"
	"github.com/gomscourse/chat-server/internal/helpers"
)

func (s *chatService) SendMessage(ctx context.Context, text string, chatID int64) error {
	sender, err := helpers.GetCtxUser(ctx)
	if err != nil {
		return err
	}

	msg, err := s.repo.CreateMessage(ctx, chatID, sender, text)
	if err != nil {
		return err
	}

	chatChan := s.initMessagesChan(chatID)

	go func() {
		chatChan <- msg
	}()

	return nil
}
