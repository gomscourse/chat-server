package chat

import (
	"fmt"
	"github.com/gomscourse/chat-server/internal/context_keys"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
	"github.com/gomscourse/chat-server/internal/service"
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
)

func (s *chatService) ConnectChat(stream service.Stream, chatID int64) error {
	ctx := stream.Context()
	username, ok := ctx.Value(context_keys.UsernameKey).(string)
	if !ok || len(username) == 0 {
		return sys.NewCommonError("invalid username in context", codes.Internal)
	}

	// проверить есть ли чат в базе и состоит ли пользователь в чате
	exists, err := s.repo.CheckUserChat(stream.Context(), chatID, username)
	if err != nil {
		return err
	}
	// если чата нет, либо пользователь не в чате - вернуть ошибку
	if !exists {
		return sys.NewCommonError("chat not found or user is not member of chat", codes.InvalidArgument)
	}

	s.mxChannel.RLock()
	chatChan, ok := s.channels[chatID]
	fmt.Println(chatChan)
	s.mxChannel.RUnlock()

	if !ok {
		s.mxChannel.Lock()
		s.channels[chatID] = make(chan *serviceModel.ChatMessage)
		s.mxChannel.Unlock()
	}

	s.mxChat.RLock()
	if _, okChat := s.chats[chatID]; !okChat {

	}
	s.mxChat.RUnlock()

	return nil
}
