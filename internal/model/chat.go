package model

import (
	"database/sql"
	"time"
)

type ChatMessage struct {
	ID        int64        `db:"id"`
	ChatID    int64        `db:"chat_id"`
	Author    string       `db:"author"`
	Content   string       `db:"content"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type Chat struct {
	ID        int64        `db:"id"`
	Title     string       `db:"title"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
