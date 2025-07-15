package chat

import "time"

type CreateIn struct {
	UserID   int64
	Message  string
	CreateAt time.Time
}
