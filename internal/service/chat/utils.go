package chat

import (
	"context"
	"fmt"
	"github.com/gomscourse/chat-server/internal/helpers"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
)

func (s *chatService) initMessagesChan(chatID int64) chan *serviceModel.ChatMessage {
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

func (s *chatService) checkChatAvailability(ctx context.Context, chatID int64, username string) error {
	// проверить есть ли чат в базе и состоит ли пользователь в чате
	exists, err := s.repo.CheckUserChat(ctx, chatID, username)
	if err != nil {
		return err
	}
	// если чата нет, либо пользователь не в чате - вернуть ошибку
	if !exists {
		return sys.NewCommonError("chat not found or user is not member of chat", codes.PermissionDenied)
	}
	return nil
}

func (s *chatService) checkUserChatAvailability(ctx context.Context, chatID int64) error {
	username, err := helpers.GetCtxUser(ctx)
	if err != nil {
		return err
	}

	return s.checkChatAvailability(ctx, chatID, username)
}
