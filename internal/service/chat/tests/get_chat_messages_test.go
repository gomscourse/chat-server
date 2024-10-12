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

func TestGetChatMessages(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx                context.Context
		id, page, pageSize int64
	}

	var (
		username = gofakeit.Name()
		ctx      = context.WithValue(context.Background(), context_keys.UsernameKey, username)
		mc       = minimock.NewController(t)
		chatID   = gofakeit.Int64()
		page     = gofakeit.Int64()
		pageSize = gofakeit.Int64()
		argsObj  = args{
			ctx:      ctx,
			id:       chatID,
			page:     page,
			pageSize: pageSize,
		}
		messages = []*model.ChatMessage{
			{ID: gofakeit.Int64(), ChatID: gofakeit.Int64(), Author: gofakeit.Name(), Content: gofakeit.Email()},
		}
		messagesCount = uint64(len(messages))

		checkUserError = fmt.Errorf("check user error")
		messagesError  = fmt.Errorf("repo messages error")
		countError     = fmt.Errorf("repo count error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		err                error
		want0              []*model.ChatMessage
		want1              uint64
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name:  "success case",
			args:  argsObj,
			err:   nil,
			want0: messages,
			want1: messagesCount,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CheckUserChatMock.Expect(ctx, chatID, username).Return(true, nil)
				mock.GetChatMessagesMock.Expect(ctx, chatID, page, pageSize).Return(messages, nil)
				mock.GetChatMessagesCountMock.Expect(ctx, chatID).Return(messagesCount, nil)
				return mock
			},
		},
		{
			name:  "check user error case",
			args:  argsObj,
			err:   checkUserError,
			want0: nil,
			want1: 0,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CheckUserChatMock.Expect(ctx, chatID, username).Return(false, checkUserError)
				return mock
			},
		},
		{
			name:  "check user validation case",
			args:  argsObj,
			err:   chatService.UserNotInChatOrChatNotFoundError,
			want0: nil,
			want1: 0,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CheckUserChatMock.Expect(ctx, chatID, username).Return(false, nil)
				return mock
			},
		},
		{
			name:  "repo messages error case",
			args:  argsObj,
			err:   messagesError,
			want0: nil,
			want1: 0,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CheckUserChatMock.Expect(ctx, chatID, username).Return(true, nil)
				mock.CheckUserChatMock.Expect(ctx, chatID, username).Return(true, nil)
				mock.GetChatMessagesMock.Expect(ctx, chatID, page, pageSize).Return(nil, messagesError)
				mock.GetChatMessagesCountMock.Expect(ctx, chatID).Return(messagesCount, nil)
				return mock
			},
		},
		{
			name:  "repo count error case",
			args:  argsObj,
			err:   countError,
			want0: nil,
			want1: 0,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CheckUserChatMock.Expect(ctx, chatID, username).Return(true, nil)
				mock.GetChatMessagesMock.Expect(ctx, chatID, page, pageSize).Return(messages, nil)
				mock.GetChatMessagesCountMock.Expect(ctx, chatID).Return(0, countError)
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

				result0, result1, err := service.GetChatMessagesAndCount(
					tt.args.ctx,
					tt.args.id,
					tt.args.page,
					tt.args.pageSize,
				)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want0, result0)
				require.Equal(t, tt.want1, result1)
			},
		)
	}
}
