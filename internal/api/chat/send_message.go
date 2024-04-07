package chat

import (
	"context"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	// TODO: add chatID param
	err := i.chatService.SendMessage(ctx, req.GetFrom(), req.GetText(), 1)
	if err != nil {
		return &emptypb.Empty{}, errors.Wrap(err, "failed to send message")
	}

	return &emptypb.Empty{}, nil
}
