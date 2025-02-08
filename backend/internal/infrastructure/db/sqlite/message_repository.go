package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/fhiroki/chat/internal/domain/message"
)

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) message.Repository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, msg *message.Message) error {
	query := `
		INSERT INTO messages (user_id, content, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		msg.UserID,
		msg.Content,
		msg.CreatedAt,
		msg.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	msg.ID = int(id)
	return nil
}

func (r *messageRepository) FindAll(ctx context.Context) ([]*message.Message, error) {
	query := `
		SELECT id, user_id, content, created_at, updated_at
		FROM messages
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*message.Message
	for rows.Next() {
		var msg message.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.UserID,
			&msg.Content,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, rows.Err()
}

func (r *messageRepository) FindByID(ctx context.Context, id int) (*message.Message, error) {
	query := `
		SELECT id, user_id, content, created_at, updated_at
		FROM messages
		WHERE id = ?
	`
	var msg message.Message
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&msg.ID,
		&msg.UserID,
		&msg.Content,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, message.ErrMessageNotFound
	}
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *messageRepository) Update(ctx context.Context, msg *message.Message) error {
	msg.UpdatedAt = time.Now()
	query := `
		UPDATE messages
		SET content = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query,
		msg.Content,
		msg.UpdatedAt,
		msg.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return message.ErrMessageNotFound
	}
	return nil
}

func (r *messageRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM messages WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return message.ErrMessageNotFound
	}
	return nil
}
