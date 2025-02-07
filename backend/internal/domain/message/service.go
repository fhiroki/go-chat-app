package message

import (
	"context"

	"github.com/fhiroki/chat/internal/infrastructure/db"
)

type MessageService interface {
	Create(ctx context.Context, msg *Message) error
	FindAll(ctx context.Context) ([]*Message, error)
	FindByID(ctx context.Context, id int) (*Message, error)
	Update(ctx context.Context, msg *Message) error
	Delete(ctx context.Context, id int) error
}

type messageService struct {
	repository Repository
	txManager  db.TxManager
}

func NewMessageService(repository Repository, txManager db.TxManager) MessageService {
	return &messageService{
		repository: repository,
		txManager:  txManager,
	}
}

func (s *messageService) Create(ctx context.Context, msg *Message) error {
	if err := msg.Validate(); err != nil {
		return ErrInvalidMessage
	}

	return s.txManager.RunInTx(ctx, func(ctx context.Context) error {
		return s.repository.Create(ctx, msg)
	})
}

func (s *messageService) FindAll(ctx context.Context) ([]*Message, error) {
	return s.repository.FindAll(ctx)
}

func (s *messageService) FindByID(ctx context.Context, id int) (*Message, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *messageService) Update(ctx context.Context, msg *Message) error {
	if err := msg.Validate(); err != nil {
		return ErrInvalidMessage
	}

	return s.txManager.RunInTx(ctx, func(ctx context.Context) error {
		return s.repository.Update(ctx, msg)
	})
}

func (s *messageService) Delete(ctx context.Context, id int) error {
	return s.txManager.RunInTx(ctx, func(ctx context.Context) error {
		return s.repository.Delete(ctx, id)
	})
}
