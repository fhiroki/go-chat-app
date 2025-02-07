package user

import (
	"time"

	"github.com/markbates/goth"
)

type User struct {
	ID        string
	Name      string
	Email     string
	AvatarURL string
	CreatedAt time.Time
}

type UserInterface interface {
	UniqueID() string
	AvatarURL() string
}

type ChatUser struct {
	goth.User
	uniqueID string
}

func (u *ChatUser) UniqueID() string {
	return u.uniqueID
}

func (u *ChatUser) SetUniqueID(id string) {
	u.uniqueID = id
}

func (u *ChatUser) AvatarURL() string {
	return u.User.AvatarURL
}
