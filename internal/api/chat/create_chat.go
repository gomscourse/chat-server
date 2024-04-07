package chat

import (
	"context"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
	"github.com/pkg/errors"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	usernames := req.GetUsernames()

	id, err := i.chatService.CreateChat(ctx, usernames)
	if err != nil {
		return &desc.CreateResponse{}, errors.Wrap(err, "failed to create chat")
	}

	return &desc.CreateResponse{Id: id}, nil
}
