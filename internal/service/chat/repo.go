package chat

import (
	"context"
	"database/sql"
	"time"
)

type Chat struct {
	ID        int64
	UserID    int64
	Message   string
	CreatedAt time.Time
}
type ChatRepository interface {
	Insert(ctx context.Context, c Chat) error
	FindByUserID(ctx context.Context, userID int64) ([]Chat, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) ChatRepository {
	return &repository{db: db}
}

func (r *repository) Insert(ctx context.Context, c Chat) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO chat (user_id, message, created_at) VALUES ($1, $2, $3)",
		c.UserID, c.Message, c.CreatedAt,
	)
	return err
}

func (r *repository) FindByUserID(ctx context.Context, userID int64) ([]Chat, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, user_id, message, created_at FROM chat WHERE user_id = $1 ORDER BY created_at ASC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var c Chat
		if err := rows.Scan(&c.ID, &c.UserID, &c.Message, &c.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, c)
	}
	return chats, nil
}
