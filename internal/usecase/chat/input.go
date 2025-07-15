package chat

type CreateIn struct {
	UserId  int64  `json:"user_id"`
	Message string `json:"message"`
}
