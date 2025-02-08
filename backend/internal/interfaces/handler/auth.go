package handler

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fhiroki/chat/internal/domain/avatar"
	"github.com/fhiroki/chat/internal/domain/user"
	userDomain "github.com/fhiroki/chat/internal/domain/user"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	userService user.UserService
}

func NewAuthHandler(userService user.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func getUserID(email string) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(email))
	return fmt.Sprintf("%x", m.Sum(nil))
}

// LoginHandler handles auth paths for Gin
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	segs := strings.Split(c.Request.URL.Path, "/")
	if len(segs) < 4 {
		c.String(404, "Auth path not found")
		return
	}

	action := segs[2]
	provider := segs[3]
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return provider, nil
	}

	switch action {
	case "login":
		gothic.BeginAuthHandler(c.Writer, c.Request)
	case "callback":
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		userID := getUserID(user.Email)
		chatUser := &userDomain.ChatUser{User: user}
		chatUser.SetUniqueID(userID)
		avatarURL, err := avatar.Avatars.GetAvatarURL(chatUser)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		userData := map[string]string{
			"user_id":    userID,
			"name":       user.Name,
			"email":      user.Email,
			"avatar_url": avatarURL,
		}

		userJSON, _ := json.Marshal(userData)
		encoded := base64.StdEncoding.EncodeToString(userJSON)

		c.SetCookie("auth", encoded, int(time.Hour*24*7/time.Second), "/", c.Request.Host, false, true)

		// ユーザー情報の保存
		u := userDomain.User{
			ID:        userID,
			Name:      user.Name,
			Email:     user.Email,
			AvatarURL: avatarURL,
			CreatedAt: time.Now(),
		}
		if err := h.userService.Create(c.Request.Context(), u); err != nil {
			c.String(500, err.Error())
			return
		}

		c.Redirect(302, "/chat")
	}
}

// LogoutHandler handles logout for Gin
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	c.SetCookie("auth", "", -1, "/", c.Request.Host, false, true)
	c.Redirect(302, "/login")
}
