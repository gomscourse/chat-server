package tests

import (
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/chat-server/internal/repository"
)

type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
