package chat

import (
	"context"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
)

func (s *chatService) GetChatMessagesAndCount(
	ctx context.Context,
	chatID, page, pageSize int64,
) ([]*serviceModel.ChatMessage, uint64, error) {
	err := s.CheckCtxUserChatAvailability(ctx, chatID)
	if err != nil {
		return nil, 0, err
	}
	var messages []*serviceModel.ChatMessage
	var count uint64
	errChan := make(chan error, 2)
	errSlice := make([]error, 0, 2)
	defer close(errChan)

	go func() {
		var err error
		messages, err = s.repo.GetChatMessages(ctx, chatID, page, pageSize)
		errChan <- err
	}()

	go func() {
		var err error
		count, err = s.repo.GetChatMessagesCount(ctx, chatID)
		errChan <- err
	}()

	for idx := 0; idx < cap(errChan); idx++ {
		if err := <-errChan; err != nil {
			errSlice = append(errSlice, err)
		}
	}

	if len(errSlice) > 0 {
		return nil, 0, errSlice[0]
	}

	return messages, count, nil
}
