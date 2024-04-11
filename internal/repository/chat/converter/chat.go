package converter

import (
	serviceModel "github.com/gomscourse/chat-server/internal/model"
	repoModel "github.com/gomscourse/chat-server/internal/repository/chat/model"
)

func ToChatMessageFromRepo(message *repoModel.ChatMessage) *serviceModel.ChatMessage {
	return &serviceModel.ChatMessage{
		ID:        message.ID,
		ChatID:    message.ChatID,
		Author:    message.Author,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}

func ToChatMessagesFromRepo(messages []*repoModel.ChatMessage) []*serviceModel.ChatMessage {
	result := make([]*serviceModel.ChatMessage, 0, len(messages))

	for _, m := range messages {
		result = append(result, ToChatMessageFromRepo(m))
	}

	return result
}
