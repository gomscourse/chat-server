package helpers

import (
	"context"
	"github.com/gomscourse/chat-server/internal/context_keys"
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
)

func GetCtxUser(ctx context.Context) (string, error) {
	username, ok := ctx.Value(context_keys.UsernameKey).(string)
	if !ok || len(username) == 0 {
		return "", sys.NewCommonError("invalid username in context", codes.Internal)
	}

	return username, nil
}
