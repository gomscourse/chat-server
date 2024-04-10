package chat

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/gomscourse/chat-server/internal/client/db"
	"github.com/gomscourse/chat-server/internal/repository"
	"github.com/pkg/errors"
)

type repo struct {
	db db.Client
}

func NewChatRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r repo) CreateChat(ctx context.Context) (int64, error) {
	var chatId int64

	q := db.Query{
		Name:     "create_chat_query",
		QueryRow: "INSERT INTO chat DEFAULT VALUES RETURNING id",
	}

	err := r.db.DB().QueryRowContext(ctx, q).Scan(&chatId)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert chat")
	}

	return chatId, nil
}

func (r repo) DeleteChat(ctx context.Context, id int64) error {
	deleteBuilder := sq.Delete("chat").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": id})
	query, args, err := deleteBuilder.ToSql()

	q := db.Query{
		Name:     "delete_chat_query",
		QueryRow: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
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

	q := db.Query{
		Name:     "add_users_to_chat_query",
		QueryRow: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
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

	q := db.Query{
		Name:     "create_message_query",
		QueryRow: query,
	}

	var messageId int64

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&messageId)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create message")
	}

	return messageId, nil
}
