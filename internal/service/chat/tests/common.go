package tests

import (
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/chat-server/internal/repository"
	"github.com/gomscourse/chat-server/internal/service"
)

type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService
type userClientMockFunc func(mc *minimock.Controller) service.UserClient
