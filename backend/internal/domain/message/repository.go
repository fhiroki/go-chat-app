package message

import (
	"context"
	"errors"
)

// ドメインエラーの定義
var (
	ErrMessageNotFound = errors.New("message not found")
	ErrInvalidMessage  = errors.New("invalid message")
)

// Repository はメッセージの永続化を担当するインターフェース
type Repository interface {
	Create(ctx context.Context, msg *Message) error
	FindAll(ctx context.Context) ([]*Message, error)
	FindByID(ctx context.Context, id int) (*Message, error)
	Update(ctx context.Context, msg *Message) error
	Delete(ctx context.Context, id int) error
}
