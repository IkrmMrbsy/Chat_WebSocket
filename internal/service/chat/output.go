package chat

import "time"

type ChatHistoryOut struct {
	ID        int64
	UserID    int64
	Message   string
	CreatedAt time.Time
}
