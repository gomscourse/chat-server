package chat

import (
	"fmt"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
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
