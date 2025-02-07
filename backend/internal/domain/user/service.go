package user

import (
	"context"

	"github.com/fhiroki/chat/internal/infrastructure/db"
)

type UserService interface {
	Create(ctx context.Context, user User) error
	FindAll(ctx context.Context) ([]*User, error)
	FindByID(ctx context.Context, id int) (*User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id int) error
}

type userService struct {
	repository Repository
	txManager  db.TxManager
}

func NewUserService(repository Repository, txManager db.TxManager) UserService {
	return &userService{repository: repository, txManager: txManager}
}

func (s *userService) Create(ctx context.Context, user User) error {
	return s.txManager.RunInTx(ctx, func(ctx context.Context) error {
		return s.repository.Create(ctx, user)
	})
}

func (s *userService) FindAll(ctx context.Context) ([]*User, error) {
	return s.repository.FindAll(ctx)
}

func (s *userService) FindByID(ctx context.Context, id int) (*User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, user User) error {
	return s.txManager.RunInTx(ctx, func(ctx context.Context) error {
		return s.repository.Update(ctx, user)
	})
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.txManager.RunInTx(ctx, func(ctx context.Context) error {
		return s.repository.Delete(ctx, id)
	})
}
