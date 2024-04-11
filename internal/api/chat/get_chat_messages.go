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
	errs := make(chan error, 2)
	defer close(errs)

	go func() {
		var err error
		messages, err = i.chatService.GetChatMessages(ctx, req.GetId(), req.GetPage(), req.GetPageSize())
		errs <- err
	}()

	go func() {
		var err error
		count, err = i.chatService.GetChatMessagesCount(ctx, req.GetId())
		errs <- err
	}()

	for idx := 0; idx < cap(errs); idx++ {
		err := <-errs
		if err != nil {
			return &desc.GetChatMessagesResponse{}, err
		}
	}

	return &desc.GetChatMessagesResponse{
		Messages: converter.ToChatMessagesFromService(messages),
		Count:    count,
	}, nil
}
