package chat

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/gomscourse/chat-server/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type repo struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) repository.ChatRepository {
	return &repo{db: db}
}

func (r repo) CreateChat(ctx context.Context) (int64, error) {
	var chatId int64
	err := r.db.QueryRow(ctx, "INSERT INTO chat DEFAULT VALUES RETURNING id").Scan(&chatId)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert chat")
	}

	return chatId, nil
}

func (r repo) DeleteChat(ctx context.Context, id int64) error {
	deleteBuilder := sq.Delete("chat").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": id})
	query, args, err := deleteBuilder.ToSql()

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete chat")
	}

	return nil
}

func (r repo) AddUsersToChat(ctx context.Context, chatID int64, usernames []string) error {
	builderInsertUserChat := sq.Insert("user_chat").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "username")

	for _, username := range usernames {
		builderInsertUserChat = builderInsertUserChat.Values(chatID, username)
	}

	query, args, err := builderInsertUserChat.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build chat query")
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to create user chat")
	}

	return nil
}

func (r repo) CreateMessage(ctx context.Context, chatID int64, sender string, text string) (int64, error) {
	builderInsertMessage := sq.Insert("message").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "author", "content").
		Values(chatID, sender, text).
		Suffix("RETURNING id")

	query, args, err := builderInsertMessage.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build message query")
	}

	var messageId int64

	err = r.db.QueryRow(ctx, query, args...).Scan(&messageId)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create message")
	}

	return messageId, nil
}
