package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/chat-server/internal/api/chat"
	"github.com/gomscourse/chat-server/internal/model"
	"github.com/gomscourse/chat-server/internal/service"
	serviceMocks "github.com/gomscourse/chat-server/internal/service/mocks"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestGetChatMessages(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *desc.GetChatMessagesRequest
	}

	var apiUpdated *timestamppb.Timestamp

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id            = gofakeit.Int64()
		page          = gofakeit.Int64()
		pageSize      = gofakeit.Int64()
		messagesCount = gofakeit.Uint64()

		serviceError = fmt.Errorf("service error")

		messageID      = gofakeit.Int64()
		messageChatID  = gofakeit.Int64()
		messageAuthor  = gofakeit.Name()
		messageContent = gofakeit.Email()
		messageUpdated = sql.NullTime{}
		messageCreated = time.Now()

		modelMessage = &model.ChatMessage{
			ID:        messageID,
			ChatID:    messageChatID,
			Author:    messageAuthor,
			Content:   messageContent,
			UpdatedAt: messageUpdated,
			CreatedAt: messageCreated,
		}
		serviceMessages = []*model.ChatMessage{modelMessage}

		apiMessage = &desc.ChatMessage{
			ID:      messageID,
			ChatID:  messageChatID,
			Author:  messageAuthor,
			Content: messageContent,
			Updated: apiUpdated,
			Created: timestamppb.New(messageCreated),
		}
		apiMessages = []*desc.ChatMessage{apiMessage}

		req = &desc.GetChatMessagesRequest{
			Id:       id,
			Page:     page,
			PageSize: pageSize,
		}

		res = &desc.GetChatMessagesResponse{
			Messages: apiMessages,
			Count:    messagesCount,
		}
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetChatMessagesResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(t)
				mock.GetChatMessagesAndCountMock.Expect(ctx, id, page, pageSize).Return(
					serviceMessages,
					messagesCount,
					nil,
				)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceError,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(t)
				mock.GetChatMessagesAndCountMock.Expect(ctx, id, page, pageSize).Return(nil, 0, serviceError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				chatServiceMock := tt.chatServiceMock(mc)
				api := chat.NewImplementation(chatServiceMock)

				result, err := api.GetChatMessages(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, result)
			},
		)
	}
}
