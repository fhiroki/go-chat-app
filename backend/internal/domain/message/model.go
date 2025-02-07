package message

import "time"

type Message struct {
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`

	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Message) AttachUserData(userData map[string]string) {
	m.UserName = userData["name"]
	m.Email = userData["email"]
	m.AvatarURL = userData["avatar_url"]
}

func (m *Message) Validate() error {
	if m.UserID == "" {
		return ErrInvalidMessage
	}
	if m.Content == "" {
		return ErrInvalidMessage
	}
	return nil
}

func NewMessage(userID, content string) *Message {
	now := time.Now()
	return &Message{
		UserID:    userID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
