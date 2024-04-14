package chat

import (
	"context"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatService.DeleteChat(ctx, req.GetId())
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
