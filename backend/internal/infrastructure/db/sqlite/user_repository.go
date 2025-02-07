package sqlite

import (
	"context"
	"database/sql"

	"github.com/fhiroki/chat/internal/domain/user"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user user.User) error {
	stmt := `INSERT OR IGNORE INTO users (id, name, email, avatar_url, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(stmt, user.ID, user.Name, user.Email, user.AvatarURL, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]*user.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, email, avatar_url, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var user user.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.AvatarURL, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*user.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, email, avatar_url, created_at FROM users WHERE id = ?", id)
	var user user.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.AvatarURL, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user user.User) error {
	stmt := `UPDATE users SET name = ?, email = ?, avatar_url = ? WHERE id = ?`
	_, err := r.db.Exec(stmt, user.Name, user.Email, user.AvatarURL, user.ID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	stmt := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(stmt, id)
	return err
}
