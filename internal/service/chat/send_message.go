package chat

import (
	"context"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
)

func (s *chatService) SendMessage(ctx context.Context, sender, text string, chatID int64) error {
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
		chatChan = make(chan *serviceModel.ChatMessage)
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
