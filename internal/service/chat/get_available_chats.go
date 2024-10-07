package chat

import (
	"context"
	"github.com/gomscourse/chat-server/internal/helpers"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
)

func (s *chatService) GetAvailableChatsAndCount(
	ctx context.Context, page, pageSize int64,
) ([]*serviceModel.Chat, uint64, error) {
	username, err := helpers.GetCtxUser(ctx)
	if err != nil {
		return nil, 0, err
	}

	var chats []*serviceModel.Chat
	var count uint64
	errChan := make(chan error, 2)
	errSlice := make([]error, 0, 2)
	defer close(errChan)

	go func() {
		var err error
		chats, err = s.GetChats(ctx, username, page, pageSize)
		errChan <- err
	}()

	go func() {
		var err error
		count, err = s.GetChatsCount(ctx, username)
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

	return chats, count, nil
}

func (s *chatService) GetChats(ctx context.Context, username string, page, pageSize int64) (
	[]*serviceModel.Chat,
	error,
) {
	return s.repo.GetChats(ctx, username, page, pageSize)
}

func (s *chatService) GetChatsCount(ctx context.Context, username string) (uint64, error) {
	return s.repo.GetChatsCount(ctx, username)
}
