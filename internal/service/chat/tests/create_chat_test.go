package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/chat-server/internal/repository"
	repositoryMocks "github.com/gomscourse/chat-server/internal/repository/mocks"
	"github.com/gomscourse/chat-server/internal/service"
	chatService "github.com/gomscourse/chat-server/internal/service/chat"
	serviceMocks "github.com/gomscourse/chat-server/internal/service/mocks"
	"github.com/gomscourse/common/pkg/db"
	commonMocks "github.com/gomscourse/common/pkg/db/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateChat(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		usernames []string
		title     string
	}

	txManagerMock := commonMocks.NewTxManagerMock(t)
	txManagerMock.ReadCommittedMock.Set(
		func(ctx context.Context, handler db.Handler) (err error) {
			return handler(ctx)
		},
	)

	var (
		ctx       = context.Background()
		mc        = minimock.NewController(t)
		usernames = []string{gofakeit.Name(), gofakeit.Name(), gofakeit.Name()}
		title     = gofakeit.Name()
		id        = gofakeit.Int64()

		checkUsersError = fmt.Errorf("repo error create")
		repoErrorCreate = fmt.Errorf("repo error create")
		repoErrorAdd    = fmt.Errorf("repo error add")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		chatRepositoryMock chatRepositoryMockFunc
		userClientMock     userClientMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:       ctx,
				usernames: usernames,
				title:     title,
			},
			want: id,
			err:  nil,
			userClientMock: func(mc *minimock.Controller) service.UserClient {
				mock := serviceMocks.NewUserClientMock(t)
				mock.CheckUsersExistenceMock.Expect(ctx, usernames).Return(nil)
				return mock
			},
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CreateChatMock.Expect(ctx, title).Return(id, nil)
				mock.AddUsersToChatMock.Expect(ctx, id, usernames).Return(nil)
				return mock
			},
		},
		{
			name: "check users error case",
			args: args{
				ctx:       ctx,
				usernames: usernames,
				title:     title,
			},
			want: 0,
			err:  checkUsersError,
			userClientMock: func(mc *minimock.Controller) service.UserClient {
				mock := serviceMocks.NewUserClientMock(t)
				mock.CheckUsersExistenceMock.Expect(ctx, usernames).Return(checkUsersError)
				return mock
			},
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CreateChatMock.Expect(ctx, title).Return(id, nil)
				mock.AddUsersToChatMock.Expect(ctx, id, usernames).Return(nil)
				return mock
			},
		},
		{
			name: "repo create error case",
			args: args{
				ctx:       ctx,
				usernames: usernames,
				title:     title,
			},
			want: 0,
			err:  repoErrorCreate,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CreateChatMock.Expect(ctx, title).Return(0, repoErrorCreate)
				return mock
			},
		},
		{
			name: "repo add error case",
			args: args{
				ctx:       ctx,
				usernames: usernames,
				title:     title,
			},
			want: 0,
			err:  repoErrorAdd,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CreateChatMock.Expect(ctx, title).Return(id, nil)
				mock.AddUsersToChatMock.Expect(ctx, id, usernames).Return(repoErrorAdd)
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
				userClientMock := tt.userClientMock(mc)
				srv := chatService.NewChatService(chatRepoMock, txManagerMock, userClientMock)

				result, err := srv.CreateChat(tt.args.ctx, tt.args.usernames, tt.args.title)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, result)
			},
		)
	}
}
