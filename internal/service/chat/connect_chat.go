package chat

import (
	"github.com/gomscourse/chat-server/internal/helpers"
	"github.com/gomscourse/chat-server/internal/logger"
	"github.com/gomscourse/chat-server/internal/service"
)

func (s *chatService) ConnectChat(stream service.Stream, chatID int64) error {
	ctx := stream.Context()
	username, err := helpers.GetCtxUser(ctx)
	if err != nil {
		return err
	}

	err = s.checkChatAvailability(ctx, chatID, username)
	if err != nil {
		return err
	}

	chatChan := s.initMessagesChan(chatID)

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
