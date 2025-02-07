package avatar

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fhiroki/chat/internal/domain/user"
	"github.com/markbates/goth"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := goth.User{
		AvatarURL: "",
	}
	testChatUser := &user.ChatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}

	testUrl := "http://url-to-avatar/"
	testChatUser.User.AvatarURL = testUrl
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	} else if url != testUrl {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &user.ChatUser{}
	user.SetUniqueID("abc")
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "https://www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.jpg")
	os.WriteFile(filename, []byte{}, 0777)
	defer os.Remove(filename)

	var fileSystemAvatar FileSystemAvatar
	user := &user.ChatUser{}
	user.SetUniqueID("abc")
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
