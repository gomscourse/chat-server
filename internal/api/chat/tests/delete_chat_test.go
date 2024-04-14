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
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestDeleteChat(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)
		id  = gofakeit.Int64()

		serviceError = fmt.Errorf("service delete error")

		req = &desc.DeleteRequest{
			Id: id,
		}

		res = &emptypb.Empty{}
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
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
				mock.DeleteChatMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  serviceError,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(t)
				mock.DeleteChatMock.Expect(ctx, id).Return(serviceError)
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

				result, err := api.Delete(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, result)
			},
		)
	}
}
