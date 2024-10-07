package chat

import (
	"context"
	"github.com/gomscourse/chat-server/internal/converter"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
)

func (i *Implementation) GetAvailableChats(
	ctx context.Context,
	req *desc.GetAvailableChatsRequest,
) (*desc.GetAvailableChatsResponse, error) {
	chats, count, err := i.chatService.GetAvailableChatsAndCount(ctx, req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, err
	}

	return &desc.GetAvailableChatsResponse{
		Chats: converter.ToChatsFromService(chats),
		Count: count,
	}, nil
}
