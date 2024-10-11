package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/chat-server/internal/context_keys"
	"github.com/gomscourse/chat-server/internal/model"
	"github.com/gomscourse/chat-server/internal/repository"
	repositoryMocks "github.com/gomscourse/chat-server/internal/repository/mocks"
	chatService "github.com/gomscourse/chat-server/internal/service/chat"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAvailableChatsAndCount(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx            context.Context
		page, pageSize int64
	}

	var (
		username = gofakeit.Name()
		ctx      = context.WithValue(context.Background(), context_keys.UsernameKey, username)
		mc       = minimock.NewController(t)
		page     = gofakeit.Int64()
		pageSize = gofakeit.Int64()
		argsObj  = args{
			ctx:      ctx,
			page:     page,
			pageSize: pageSize,
		}
		chats = []*model.Chat{
			{ID: gofakeit.Int64(), Title: gofakeit.Name()},
		}
		chatsCount = uint64(len(chats))

		chatsError      = fmt.Errorf("repo chats error")
		chatsCountError = fmt.Errorf("repo chats count error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		err                error
		allErrors          []error
		want0              []*model.Chat
		want1              uint64
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name:  "success case",
			args:  argsObj,
			err:   nil,
			want0: chats,
			want1: chatsCount,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.GetChatsMock.Expect(ctx, username, page, pageSize).Return(chats, nil)
				mock.GetChatsCountMock.Expect(ctx, username).Return(chatsCount, nil)
				return mock
			},
		},
		{
			name:  "repo chats error case",
			args:  argsObj,
			err:   chatsError,
			want0: nil,
			want1: 0,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.GetChatsMock.Expect(ctx, username, page, pageSize).Return(nil, chatsError)
				mock.GetChatsCountMock.Expect(ctx, username).Return(chatsCount, nil)
				return mock
			},
		},
		{
			name:  "repo chats count error case",
			args:  argsObj,
			err:   chatsCountError,
			want0: nil,
			want1: 0,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.GetChatsMock.Expect(ctx, username, page, pageSize).Return(chats, nil)
				mock.GetChatsCountMock.Expect(ctx, username).Return(0, chatsCountError)
				return mock
			},
		},
		{
			name:      "repo chats and count error case",
			args:      argsObj,
			err:       chatsError,
			allErrors: []error{chatsError, chatsCountError},
			want0:     nil,
			want1:     0,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.GetChatsMock.Expect(ctx, username, page, pageSize).Return(nil, chatsError)
				mock.GetChatsCountMock.Expect(ctx, username).Return(0, chatsCountError)
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

				result0, result1, err := service.GetAvailableChatsAndCount(tt.args.ctx, tt.args.page, tt.args.pageSize)
				if tt.name == "repo chats and count error case" {
					require.Contains(t, tt.allErrors, err)
				} else {
					require.Equal(t, tt.err, err)
				}
				require.Equal(t, tt.want0, result0)
				require.Equal(t, tt.want1, result1)
			},
		)
	}
}
