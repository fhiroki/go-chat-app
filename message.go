package main

import (
	"net/http"
	"strconv"
	"time"
)

type SendMessage struct {
	ID        int
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

func deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	// CORSの許可設定: 任意のオリジンからのリクエストを許可
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")

	// プリフライトリクエストの場合はここで終了
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := DeleteMessage(idInt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
	stmt := `SELECT m.id, m.message, u.name, u.email, u.avatar_url, m.created_at
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
		err := rows.Scan(&msg.ID, &msg.Message, &msg.UserName, &msg.Email, &msg.AvatarURL, &msg.CreatedAt)
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
