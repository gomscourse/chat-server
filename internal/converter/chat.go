package converter

import (
	serviceModel "github.com/gomscourse/chat-server/internal/model"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToChatMessageFromService(message *serviceModel.ChatMessage) *desc.ChatMessage {
	var updatedAt *timestamppb.Timestamp
	if message.UpdatedAt.Valid {
		updatedAt = timestamppb.New(message.UpdatedAt.Time)
	}

	return &desc.ChatMessage{
		ID:      message.ID,
		ChatID:  message.ChatID,
		Author:  message.Author,
		Content: message.Content,
		Created: timestamppb.New(message.CreatedAt),
		Updated: updatedAt,
	}
}

func ToChatMessagesFromService(messages []*serviceModel.ChatMessage) []*desc.ChatMessage {
	result := make([]*desc.ChatMessage, 0, len(messages))

	for _, m := range messages {
		result = append(result, ToChatMessageFromService(m))
	}

	return result
}
