package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/chat-server/internal/model"
	"github.com/gomscourse/chat-server/internal/repository"
	repositoryMocks "github.com/gomscourse/chat-server/internal/repository/mocks"
	chatService "github.com/gomscourse/chat-server/internal/service/chat"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetChatMessages(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx                context.Context
		id, page, pageSize int64
	}

	var (
		ctx      = context.Background()
		mc       = minimock.NewController(t)
		id       = gofakeit.Int64()
		page     = gofakeit.Int64()
		pageSize = gofakeit.Int64()
		argsObj  = args{
			ctx:      ctx,
			id:       id,
			page:     page,
			pageSize: pageSize,
		}
		messages = []*model.ChatMessage{
			{ID: gofakeit.Int64(), ChatID: gofakeit.Int64(), Author: gofakeit.Name(), Content: gofakeit.Email()},
		}

		retrieveError = fmt.Errorf("repo retrieve error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		err                error
		want               []*model.ChatMessage
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name: "success case",
			args: argsObj,
			err:  nil,
			want: messages,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.GetChatMessagesMock.Expect(ctx, id, page, pageSize).Return(messages, nil)
				return mock
			},
		},
		{
			name: "repo error case",
			args: argsObj,
			err:  retrieveError,
			want: nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.GetChatMessagesMock.Expect(ctx, id, page, pageSize).Return(nil, retrieveError)
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

				result, err := service.GetChatMessages(tt.args.ctx, tt.args.id, tt.args.page, tt.args.pageSize)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, result)
			},
		)
	}
}

func TestGetChatMessagesCount(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		id      = gofakeit.Int64()
		count   = gofakeit.Uint64()
		argsObj = args{
			ctx: ctx,
			id:  id,
		}

		countError = fmt.Errorf("repo count error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		err                error
		want               uint64
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name: "success case",
			args: argsObj,
			err:  nil,
			want: count,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.GetChatMessagesCountMock.Expect(ctx, id).Return(count, nil)
				return mock
			},
		},
		{
			name: "repo error case",
			args: argsObj,
			err:  countError,
			want: 0,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.GetChatMessagesCountMock.Expect(ctx, id).Return(0, countError)
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

				result, err := service.GetChatMessagesCount(tt.args.ctx, tt.args.id)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, result)
			},
		)
	}
}
