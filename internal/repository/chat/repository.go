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
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const messageTableName = "message"

const (
	idColumn     = "id"
	chatIdColumn = "chat_id"
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

func (r repo) CreateChat(ctx context.Context, title string) (int64, error) {
	var chatId int64
	builderInsertMessage := sq.Insert("chat").
		PlaceholderFormat(sq.Dollar).
		Columns("title").
		Values(title).
		Suffix("RETURNING id")

	query, args, err := builderInsertMessage.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build message query")
	}

	q := db.Query{
		Name:     "create_chat_query",
		QueryRow: query,
	}

	err = r.db.DB().QueryRowContextScan(ctx, &chatId, q, args...)
	//TODO: обработать ошибку из-за существующего title
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

func (r repo) CreateMessage(ctx context.Context, chatID int64, sender string, text string) (
	*serviceModel.ChatMessage,
	error,
) {
	builderInsertMessage := sq.Insert("message").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "author", "content").
		Values(chatID, sender, text).
		Suffix("RETURNING id, created_at")

	query, args, err := builderInsertMessage.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build message query")
	}

	q := db.Query{
		Name:     "create_message_query",
		QueryRow: query,
	}

	var messageId int64
	var created pgtype.Timestamp

	err = r.db.DB().QueryRowContextScanMany(ctx, []any{&messageId, &created}, q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create message")
	}

	return &serviceModel.ChatMessage{
		ID:        messageId,
		ChatID:    chatID,
		Author:    sender,
		Content:   text,
		CreatedAt: created.Time,
	}, nil
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
		Where(sq.Eq{chatIdColumn: chatID}).
		OrderBy("id DESC").
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

	//TODO: добавить аналогичную ошибку в common и подменять на нее при возврате
	// чтобы не зависеть от пакет pgx
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, sys.NewCommonError(fmt.Sprintf("messages for chat with id %d not found", chatID), codes.NotFound)
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

func (r repo) CheckUserChat(ctx context.Context, chatID int64, username string) (bool, error) {
	q := db.Query{
		Name:     "check_user_chat",
		QueryRow: "SELECT count(id) FROM user_chat WHERE chat_id = $1 AND username = $2;",
	}

	var count int
	err := r.db.DB().QueryRowContextScan(ctx, &count, q, chatID, username)
	if err != nil {
		return false, errors.Wrap(err, "failed to get user_chat")
	}

	return count > 0, nil
}
