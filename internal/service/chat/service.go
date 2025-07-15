package chat

import (
	"context"
	"time"
)

type Service interface {
	HandleMessage(ctx context.Context, input CreateIn) error
	GetHistory(ctx context.Context, UserID int64) ([]ChatHistoryOut, error)
}

type service struct {
	repo ChatRepository
}

func NewService(repo ChatRepository) Service {
	return &service{repo}
}

func (s *service) HandleMessage(ctx context.Context, input CreateIn) error {
	chat := Chat{
		UserID:    input.UserID,
		Message:   input.Message,
		CreatedAt: time.Now(),
	}

	return s.repo.Insert(ctx, chat)
}

func (s *service) GetHistory(ctx context.Context, userID int64) ([]ChatHistoryOut, error) {
	chats, err := s.repo.FindByUserID(ctx, userID)
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
