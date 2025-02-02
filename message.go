package main

import (
	"time"
)

type SendMessage struct {
	Message   string
	UserName  string
	Email     string
	AvatarURL string
	CreatedAt time.Time
}

type Message struct {
	ID        int
	UserID    string
	Message   string
	CreatedAt time.Time
}

func SaveMessage(msg Message) (int64, error) {
	stmt := `INSERT INTO messages (user_id, message, created_at) VALUES (?, ?, ?)`
	res, err := db.Exec(stmt, msg.UserID, msg.Message, msg.CreatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetMessages() ([]SendMessage, error) {
	stmt := `SELECT m.message, u.name, u.email, u.avatar_url, m.created_at
	FROM messages m
	JOIN users u ON m.user_id = u.id
	ORDER BY m.created_at ASC LIMIT 100`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []SendMessage
	for rows.Next() {
		var msg SendMessage
		err := rows.Scan(&msg.Message, &msg.UserName, &msg.Email, &msg.AvatarURL, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func GetMessageByID(id int) (Message, error) {
	var msg Message
	row := db.QueryRow("SELECT id, user_id, message, created_at FROM messages WHERE id = ?", id)
	err := row.Scan(&msg.ID, &msg.UserID, &msg.Message, &msg.CreatedAt)
	return msg, err
}

func UpdateMessage(msg Message) error {
	stmt := `UPDATE messages SET user_id = ?, message = ? WHERE id = ?`
	_, err := db.Exec(stmt, msg.UserID, msg.Message, msg.ID)
	return err
}

func DeleteMessage(id int) error {
	stmt := `DELETE FROM messages WHERE id = ?`
	_, err := db.Exec(stmt, id)
	return err
}
