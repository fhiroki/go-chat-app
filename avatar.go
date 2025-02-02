package main

import (
	"errors"
	"fmt"
)

var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL.")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userID, ok := c.userData["user_id"]; ok {
		return fmt.Sprintf("https://www.gravatar.com/avatar/%s", userID), nil
	}
	return "", ErrNoAvatarURL
}
