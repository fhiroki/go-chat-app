package avatar

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fhiroki/chat/internal/domain/user"
)

var Avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL.")

type Avatar interface {
	GetAvatarURL(user.UserInterface) (string, error)
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u user.UserInterface) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(u user.UserInterface) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(u user.UserInterface) (string, error) {
	return fmt.Sprintf("https://www.gravatar.com/avatar/%s", u.UniqueID()), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(u user.UserInterface) (string, error) {
	if files, err := os.ReadDir("assets/avatars"); err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasPrefix(file.Name(), u.UniqueID()) {
				return "/assets/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
