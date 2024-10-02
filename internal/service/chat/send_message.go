package chat

import (
	"context"
	"github.com/gomscourse/chat-server/internal/helpers"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
)

func (s *chatService) SendMessage(ctx context.Context, text string, chatID int64) error {
	sender, err := helpers.GetCtxUser(ctx)
	if err != nil {
		return err
	}

	//TODO: переделать на возврат модели сообщения
	id, err := s.repo.CreateMessage(ctx, chatID, sender, text)
	if err != nil {
		return err
	}

	s.mxChannel.RLock()
	chatChan, ok := s.channels[chatID]
	s.mxChannel.RUnlock()

	if !ok {
		s.mxChannel.Lock()
		chatChan = make(chan *serviceModel.ChatMessage, 100)
		s.channels[chatID] = chatChan
		s.mxChannel.Unlock()
	}

	go func() {
		chatChan <- &serviceModel.ChatMessage{
			ID:      id,
			ChatID:  chatID,
			Author:  sender,
			Content: text,
		}
	}()

	return nil
}
