package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/chat-server/internal/context_keys"
	"github.com/gomscourse/chat-server/internal/model"
	"github.com/gomscourse/chat-server/internal/repository"
	repositoryMocks "github.com/gomscourse/chat-server/internal/repository/mocks"
	"github.com/gomscourse/chat-server/internal/service"
	chatService "github.com/gomscourse/chat-server/internal/service/chat"
	serviceMocks "github.com/gomscourse/chat-server/internal/service/mocks"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx          context.Context
		id           int64
		sender, text string
	}

	var (
		author  = gofakeit.Name()
		ctx     = context.WithValue(context.Background(), context_keys.UsernameKey, author)
		mc      = minimock.NewController(t)
		msgID   = gofakeit.Int64()
		chatID  = gofakeit.Int64()
		content = gofakeit.Email()

		ch = make(chan *model.ChatMessage, 100)

		sendError = fmt.Errorf("repo error delete")

		msg = &model.ChatMessage{
			ID:        msgID,
			ChatID:    chatID,
			Author:    author,
			Content:   content,
			UpdatedAt: sql.NullTime{},
			CreatedAt: time.Now(),
		}
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		err                error
		want               *model.ChatMessage
		chatRepositoryMock chatRepositoryMockFunc
		chatServiceMock    chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				id:     chatID,
				sender: author,
				text:   content,
			},
			err: nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CreateMessageMock.Expect(ctx, chatID, author, content).Return(msg, nil)
				return mock
			},
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(t)
				mock.InitMessagesChanMock.Expect(chatID).Return(ch)
				return mock
			},
		},
		{
			name: "repo error case",
			args: args{
				ctx:    ctx,
				id:     chatID,
				sender: author,
				text:   content,
			},
			err: sendError,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(t)
				mock.CreateMessageMock.Expect(ctx, chatID, author, content).Return(nil, sendError)
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
				srv := chatService.NewTestService(chatRepoMock)

				err := srv.SendMessage(tt.args.ctx, tt.args.text, tt.args.id)
				require.Equal(t, tt.err, err)
				if err != nil {
					m := <-ch
					require.Equal(t, tt.want, m)
				}
			},
		)
	}
}
