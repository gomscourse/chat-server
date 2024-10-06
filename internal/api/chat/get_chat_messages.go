package chat

import (
	"context"
	"github.com/gomscourse/chat-server/internal/converter"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
)

func (i *Implementation) GetChatMessages(
	ctx context.Context,
	req *desc.GetChatMessagesRequest,
) (*desc.GetChatMessagesResponse, error) {
	messages, count, err := i.chatService.GetChatMessagesAndCount(ctx, req.GetId(), req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, err
	}

	return &desc.GetChatMessagesResponse{
		Messages: converter.ToChatMessagesFromService(messages),
		Count:    count,
	}, nil
}
