package chat

import (
	"context"
	"wschat/internal/service/chat"
)

type ChatUsecase interface {
	HanldeIncomingMessage(ctx context.Context, input CreateIn) error
	GetHistory(ctx context.Context, userID int64) ([]ChatHistoryOut, error)
}

type chatUsecase struct {
	service chat.Service
}

func NewChatUsecase(service chat.Service) ChatUsecase {
	return &chatUsecase{service}
}

func (u *chatUsecase) HanldeIncomingMessage(ctx context.Context, input CreateIn) error {
	return u.service.HandleMessage(ctx, chat.CreateIn{
		UserID:  input.UserId,
		Message: input.Message,
	})
}

func (u *chatUsecase) GetHistory(ctx context.Context, userID int64) ([]ChatHistoryOut, error) {
	chats, err := u.service.GetHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	var out []ChatHistoryOut
	for _, c := range chats {
		out = append(out, ChatHistoryOut{
			ID:        c.ID,
			UserID:    c.UserID,
			Message:   c.Message,
			CreatedAt: c.CreatedAt,
		})
	}

	return out, nil
}
