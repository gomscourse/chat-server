package chat

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	serviceModel "github.com/gomscourse/chat-server/internal/model"
	"github.com/gomscourse/chat-server/internal/repository"
	"github.com/gomscourse/chat-server/internal/repository/chat/converter"
	repoModel "github.com/gomscourse/chat-server/internal/repository/chat/model"
	"github.com/gomscourse/common/pkg/db"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const messageTableName = "message"

const (
	idColumn = "id"
)

const (
	messageChatIDColumn  = "chat_id"
	messageAuthorColumn  = "author"
	messageContentColumn = "content"
)

const (
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
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

	err := r.db.DB().QueryRowContextScan(ctx, &chatId, q)
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

	err = r.db.DB().QueryRowContextScan(ctx, &messageId, q, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create message")
	}

	return messageId, nil
}

func (r repo) GetChatMessages(ctx context.Context, chatID, page, pageSize int64) ([]*serviceModel.ChatMessage, error) {
	limit := uint64(pageSize)
	offset := uint64((page - 1) * pageSize)
	builderSelect := sq.Select(
		idColumn,
		messageChatIDColumn,
		messageAuthorColumn,
		messageContentColumn,
		createdAtColumn,
		updatedAtColumn,
	).
		From(messageTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: chatID}).
		Limit(limit).
		Offset(offset)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "get_chat_messages_query",
		QueryRow: query,
	}

	var messages []*repoModel.ChatMessage
	err = r.db.DB().ScanAllContext(ctx, &messages, q, args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("messages for chat with id %d not found", chatID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return converter.ToChatMessagesFromRepo(messages), nil
}

func (r repo) GetChatMessagesCount(ctx context.Context, chatID int64) (uint64, error) {
	q := db.Query{
		Name:     "chat_messages_count",
		QueryRow: "SELECT COUNT(id) FROM message WHERE chat_id = $1;",
	}

	var count uint64

	err := r.db.DB().QueryRowContextScan(ctx, &count, q, chatID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count chat messages")
	}

	return count, nil
}
