package chat

import (
	"context"
	"github.com/gomscourse/chat-server/internal/converter"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
)

func (i *Implementation) GetChatMessages(ctx context.Context, req *desc.GetChatMessagesRequest) (*desc.GetChatMessagesResponse, error) {
	var messages []*serviceModel.ChatMessage
	var count uint64
	errChan := make(chan error, 2)
	errSlice := make([]error, 0, 2)
	defer close(errChan)

	go func() {
		var err error
		messages, err = i.chatService.GetChatMessages(ctx, req.GetId(), req.GetPage(), req.GetPageSize())
		errChan <- err
	}()

	go func() {
		var err error
		count, err = i.chatService.GetChatMessagesCount(ctx, req.GetId())
		errChan <- err
	}()

	for idx := 0; idx < cap(errChan); idx++ {
		if err := <-errChan; err != nil {
			errSlice = append(errSlice, err)
		}
	}

	if len(errSlice) > 0 {
		return &desc.GetChatMessagesResponse{}, errSlice[0]
	}

	return &desc.GetChatMessagesResponse{
		Messages: converter.ToChatMessagesFromService(messages),
		Count:    count,
	}, nil
}
