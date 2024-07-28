package chat

import (
	"context"
	"github.com/gomscourse/chat-server/internal/converter"
	"github.com/gomscourse/chat-server/internal/model"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
)

type streamWrapper struct {
	stream desc.ChatV1_ConnectChatServer
}

func NewStreamWrapper(stream desc.ChatV1_ConnectChatServer) *streamWrapper {
	return &streamWrapper{
		stream: stream,
	}
}

func (w *streamWrapper) Send(message *model.ChatMessage) error {
	return w.stream.Send(converter.ToChatMessageFromService(message))
}

func (w *streamWrapper) Context() context.Context {
	return w.stream.Context()
}

func (i *Implementation) ConnectChat(req *desc.ConnectChatRequest, stream desc.ChatV1_ConnectChatServer) error {
	return i.chatService.ConnectChat(NewStreamWrapper(stream), req.GetChatId())
}
