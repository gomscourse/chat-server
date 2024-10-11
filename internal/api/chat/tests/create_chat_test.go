package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/chat-server/internal/api/chat"
	"github.com/gomscourse/chat-server/internal/service"
	serviceMocks "github.com/gomscourse/chat-server/internal/service/mocks"
	desc "github.com/gomscourse/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateChat(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx       = context.Background()
		mc        = minimock.NewController(t)
		usernames = []string{gofakeit.Name(), gofakeit.Name(), gofakeit.Name()}
		chatTitle = gofakeit.Name()
		id        = gofakeit.Int64()

		serviceError = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Usernames: usernames,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
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
				mock.CreateChatMock.Expect(ctx, usernames, chatTitle).Return(id, nil)
				return mock
			},
		},
		{
			name: "service create error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceError,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(t)
				mock.CreateChatMock.Expect(ctx, usernames, chatTitle).Return(0, serviceError)
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

				result, err := api.Create(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, result)
			},
		)
	}
}
