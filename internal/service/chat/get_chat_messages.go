package chat

import (
	"context"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
)

func (s *chatService) GetChatMessages(ctx context.Context, chatID, page, pageSize int64) (
	[]*serviceModel.ChatMessage,
	error,
) {
	err := s.checkUserChatAvailability(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetChatMessages(ctx, chatID, page, pageSize)
}

func (s *chatService) GetChatMessagesCount(ctx context.Context, chatID int64) (uint64, error) {
	err := s.checkUserChatAvailability(ctx, chatID)
	if err != nil {
		return 0, err
	}

	return s.repo.GetChatMessagesCount(ctx, chatID)
}
