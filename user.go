package main

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	AvatarURL string
	CreatedAt time.Time
}

func CreateUser(user User) error {
	stmt := `INSERT OR IGNORE INTO users (id, name, email, avatar_url, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(stmt, user.ID, user.Name, user.Email, user.AvatarURL, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByID(id int) (User, error) {
	var user User
	row := db.QueryRow("SELECT id, name, email, avatar_url, created_at FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.AvatarURL, &user.CreatedAt)
	return user, err
}

func UpdateUser(user User) error {
	stmt := `UPDATE users SET name = ?, email = ?, avatar_url = ? WHERE id = ?`
	_, err := db.Exec(stmt, user.Name, user.Email, user.AvatarURL, user.ID)
	return err
}

func DeleteUser(id int) error {
	stmt := `DELETE FROM users WHERE id = ?`
	_, err := db.Exec(stmt, id)
	return err
}
