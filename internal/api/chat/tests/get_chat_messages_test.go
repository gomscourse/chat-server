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

		serviceRetrieveError = fmt.Errorf("service get messages error")
		serviceCountError    = fmt.Errorf("service get messages error")

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
				mock.GetChatMessagesMock.Expect(ctx, id, page, pageSize).Return(serviceMessages, nil)
				mock.GetChatMessagesCountMock.Expect(ctx, id).Return(messagesCount, nil)
				return mock
			},
		},
		{
			name: "service error get chat messages case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceRetrieveError,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(t)
				mock.GetChatMessagesMock.Expect(ctx, id, page, pageSize).Return(nil, serviceRetrieveError)
				mock.GetChatMessagesCountMock.Expect(ctx, id).Return(messagesCount, nil)
				return mock
			},
		},
		{
			name: "service error get chat messages count case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceCountError,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(t)
				mock.GetChatMessagesMock.Expect(ctx, id, page, pageSize).Return(serviceMessages, nil)
				mock.GetChatMessagesCountMock.Expect(ctx, id).Return(0, serviceCountError)
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
				if result != nil && len(result.Messages) > 0 {
					require.Equal(t, tt.want.Messages[0].ID, result.Messages[0].ID)
					require.Equal(t, tt.want.Messages[0].ChatID, result.Messages[0].ChatID)
					require.Equal(t, tt.want.Messages[0].Author, result.Messages[0].Author)
					require.Equal(t, tt.want.Messages[0].Content, result.Messages[0].Content)
					require.Equal(t, tt.want.Messages[0].Updated, result.Messages[0].Updated)
					require.Equal(t, tt.want.Messages[0].Created, result.Messages[0].Created)
				} else {
					require.Equal(t, tt.want, result)
				}
			},
		)
	}
}
