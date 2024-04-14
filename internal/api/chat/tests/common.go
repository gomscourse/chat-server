package tests

import (
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/chat-server/internal/service"
)

type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService
