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

const (
	messageTableName  = "message"
	chatTableName     = "chat"
	userChatTableName = "user_chat"
	usernameColumn    = "username"
)

const (
	idColumn        = "id"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

const (
	messageChatIDColumn  = "chat_id"
	messageAuthorColumn  = "author"
	messageContentColumn = "content"
)

const (
	chatTitleColumn = "title"
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
		return 0, sys.NewCommonError(errors.Wrap(err, "failed to insert chat").Error(), codes.Internal)
	}

	return chatId, nil
}

func (r repo) DeleteChat(ctx context.Context, id int64) error {
	//FIXME: также удалить записи в user_chat и сообщения (либо помечать чат как удаленный)
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
	builderSelect := prepareChatMessagesQuery(
		chatID, idColumn,
		messageChatIDColumn,
		messageAuthorColumn,
		messageContentColumn,
		createdAtColumn,
		updatedAtColumn,
	).OrderBy("id DESC")

	builderSelect = handleLimitAndOffset(builderSelect, limit, offset)

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
	builderSelect := prepareChatMessagesQuery(chatID, fmt.Sprintf("COUNT(%s)", idColumn))

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "chat_messages_count",
		QueryRow: query,
	}

	var count uint64

	err = r.db.DB().QueryRowContextScan(ctx, &count, q, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count chat messages")
	}

	return count, nil
}

func prepareChatMessagesQuery(chatID int64, selects ...string) sq.SelectBuilder {
	return sq.Select(selects...).
		From(messageTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{messageChatIDColumn: chatID})
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

func (r repo) GetChats(ctx context.Context, username string, page, pageSize int64) ([]*serviceModel.Chat, error) {
	limit := uint64(pageSize)
	offset := uint64((page - 1) * pageSize)
	builderSelect := prepareUserChatsQuery(
		username,
		fmt.Sprintf("%s.%s", chatTableName, idColumn),
		chatTitleColumn,
		fmt.Sprintf("%s.%s", chatTableName, createdAtColumn),
		fmt.Sprintf("%s.%s", chatTableName, updatedAtColumn),
	)
	builderSelect = handleLimitAndOffset(builderSelect, limit, offset)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "get_chats_query",
		QueryRow: query,
	}

	var chats []*repoModel.Chat
	err = r.db.DB().ScanAllContext(ctx, &chats, q, args...)

	//TODO: добавить аналогичную ошибку в common и подменять на нее при возврате
	// чтобы не зависеть от пакет pgx
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, sys.NewCommonError(fmt.Sprintf("chats for chat user %s not found", username), codes.NotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get chats: %w", err)
	}

	return converter.ToChatsFromRepo(chats), nil
}

func (r repo) GetChatsCount(ctx context.Context, username string) (uint64, error) {
	builderSelect := prepareUserChatsQuery(username, fmt.Sprintf("COUNT(%s.%s)", chatTableName, idColumn))

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "chats_count",
		QueryRow: query,
	}

	var count uint64

	err = r.db.DB().QueryRowContextScan(ctx, &count, q, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count chat messages")
	}

	return count, nil
}

func prepareUserChatsQuery(username string, columns ...string) sq.SelectBuilder {
	return sq.Select(columns...).
		From(chatTableName).
		PlaceholderFormat(sq.Dollar).
		InnerJoin(fmt.Sprintf("%s ON %s.chat_id = %s.id", userChatTableName, userChatTableName, chatTableName)).
		Where(sq.Eq{fmt.Sprintf("%s.%s", userChatTableName, usernameColumn): username})
}

func handleLimitAndOffset(builderSelect sq.SelectBuilder, limit, offset uint64) sq.SelectBuilder {
	if limit != 0 {
		builderSelect = builderSelect.Limit(limit)
	}

	if offset != 0 {
		builderSelect = builderSelect.Offset(offset)
	}

	return builderSelect
}
