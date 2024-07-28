package chat

import (
	"fmt"
	"github.com/gomscourse/chat-server/internal/service"
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
)

func (s *chatService) ConnectChat(stream service.Stream, chatID int64) error {
	s.mxChannel.RLock()
	chatChan, ok := s.channels[chatID]
	fmt.Println(chatChan)
	s.mxChannel.RUnlock()

	if !ok {
		return sys.NewCommonError(fmt.Sprintf("chat with id %d not found", chatID), codes.NotFound)
	}

	s.mxChat.RLock()
	if _, okChat := s.chats[chatID]; !okChat {

	}
	s.mxChat.RUnlock()

	return nil
}
