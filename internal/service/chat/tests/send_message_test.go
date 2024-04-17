package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/chat-server/internal/repository"
	repositoryMocks "github.com/gomscourse/chat-server/internal/repository/mocks"
	chatService "github.com/gomscourse/chat-server/internal/service/chat"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx          context.Context
		id           int64
		sender, text string
	}

	var (
		ctx    = context.Background()
		mc     = minimock.NewController(t)
		id     = gofakeit.Int64()
		sender = gofakeit.Name()
		text   = gofakeit.Email()

		sendError = fmt.Errorf("repo error delete")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		err                error
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				id:     id,
				sender: sender,
				text:   text,
			},
			err: nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CreateMessageMock.Expect(ctx, id, sender, text).Return(1, nil)
				return mock
			},
		},
		{
			name: "repo error case",
			args: args{
				ctx:    ctx,
				id:     id,
				sender: sender,
				text:   text,
			},
			err: sendError,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CreateMessageMock.Expect(ctx, id, sender, text).Return(0, sendError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				chatRepoMock := tt.chatRepositoryMock(mc)
				service := chatService.NewTestService(chatRepoMock)

				err := service.SendMessage(tt.args.ctx, tt.args.sender, tt.args.text, tt.args.id)
				require.Equal(t, tt.err, err)
			},
		)
	}
}
