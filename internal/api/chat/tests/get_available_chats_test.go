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

func TestGetAvailableChats(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *desc.GetAvailableChatsRequest
	}

	var apiUpdated *timestamppb.Timestamp

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		page       = gofakeit.Int64()
		pageSize   = gofakeit.Int64()
		chatsCount = gofakeit.Uint64()

		serviceRetrieveError = fmt.Errorf("service get available chats error")

		chatID         = gofakeit.Int64()
		chatTitle      = gofakeit.Name()
		messageUpdated = sql.NullTime{}
		messageCreated = time.Now()

		modelChat = &model.Chat{
			ID:        chatID,
			Title:     chatTitle,
			UpdatedAt: messageUpdated,
			CreatedAt: messageCreated,
		}
		serviceChats = []*model.Chat{modelChat}

		apiChat = &desc.Chat{
			ID:      chatID,
			Title:   chatTitle,
			Updated: apiUpdated,
			Created: timestamppb.New(messageCreated),
		}
		apiChats = []*desc.Chat{apiChat}

		req = &desc.GetAvailableChatsRequest{
			Page:     page,
			PageSize: pageSize,
		}

		res = &desc.GetAvailableChatsResponse{
			Chats: apiChats,
			Count: chatsCount,
		}
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetAvailableChatsResponse
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
				mock.GetAvailableChatsAndCountMock.Expect(ctx, page, pageSize).Return(serviceChats, chatsCount, nil)
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
				mock.GetAvailableChatsAndCountMock.Expect(ctx, page, pageSize).Return(nil, 0, serviceRetrieveError)
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

				result, err := api.GetAvailableChats(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, result)
			},
		)
	}
}
