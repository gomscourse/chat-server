package chat

import (
	"context"
	"fmt"
	"github.com/gomscourse/chat-server/internal/helpers"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
)

var UserNotInChatOrChatNotFoundError = sys.NewCommonError(
	"chat not found or user is not member of chat",
	codes.PermissionDenied,
)

func (s *chatService) InitMessagesChan(chatID int64) chan *serviceModel.ChatMessage {
	s.mxChannel.RLock()
	chatChan, ok := s.channels[chatID]
	fmt.Println(chatChan)
	s.mxChannel.RUnlock()

	if !ok {
		chatChan = make(chan *serviceModel.ChatMessage, 100)
		s.mxChannel.Lock()
		s.channels[chatID] = chatChan
		s.mxChannel.Unlock()
	}

	return chatChan
}

func (s *chatService) CheckChatAvailability(ctx context.Context, chatID int64, username string) error {
	// проверить есть ли чат в базе и состоит ли пользователь в чате
	exists, err := s.repo.CheckUserChat(ctx, chatID, username)
	if err != nil {
		return err
	}
	// если чата нет, либо пользователь не в чате - вернуть ошибку
	if !exists {
		return UserNotInChatOrChatNotFoundError
	}
	return nil
}

func (s *chatService) CheckCtxUserChatAvailability(ctx context.Context, chatID int64) error {
	username, err := helpers.GetCtxUser(ctx)
	if err != nil {
		return err
	}

	return s.CheckChatAvailability(ctx, chatID, username)
}

func (s *chatService) GetChannels() map[int64]chan *serviceModel.ChatMessage {
	return s.channels
}
