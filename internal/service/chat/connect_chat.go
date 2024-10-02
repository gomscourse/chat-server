package chat

import (
	"fmt"
	"github.com/gomscourse/chat-server/internal/helpers"
	"github.com/gomscourse/chat-server/internal/logger"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
	"github.com/gomscourse/chat-server/internal/service"
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
)

func (s *chatService) ConnectChat(stream service.Stream, chatID int64) error {
	ctx := stream.Context()
	username, err := helpers.GetCtxUser(ctx)
	if err != nil {
		return err
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
		chatChan = make(chan *serviceModel.ChatMessage, 100)
		s.mxChannel.Lock()
		s.channels[chatID] = chatChan
		s.mxChannel.Unlock()
	}

	s.mxChat.RLock()
	if _, okChat := s.chats[chatID]; !okChat {
		s.chats[chatID] = &Chat{streams: make(map[string]service.Stream)}
	}
	s.mxChat.RUnlock()

	s.chats[chatID].m.Lock()
	s.chats[chatID].streams[username] = stream
	s.chats[chatID].m.Unlock()

	for {
		select {
		case msg, okChan := <-chatChan:
			if !okChan {
				return nil
			}

			for u, st := range s.chats[chatID].streams {
				if err := st.Send(msg); err != nil {
					logger.Error(err.Error(), "chatID", chatID, "username", u)
				}
			}

		case <-ctx.Done():
			s.chats[chatID].m.Lock()
			delete(s.chats[chatID].streams, username)
			s.chats[chatID].m.Unlock()
			return nil
		}
	}
}
