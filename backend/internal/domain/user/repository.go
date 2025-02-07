package user

import "context"

type Repository interface {
	Create(ctx context.Context, user User) error
	FindAll(ctx context.Context) ([]*User, error)
	FindByID(ctx context.Context, id int) (*User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id int) error
}
