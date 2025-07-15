package user

import "context"

type User struct {
	ID   int64
	Name string
}

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
}
