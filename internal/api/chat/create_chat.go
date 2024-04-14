package chat

import (
	"context"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	usernames := req.GetUsernames()

	id, err := i.chatService.CreateChat(ctx, usernames)
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{Id: id}, nil
}
