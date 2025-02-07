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

// LoginHandler handles auth paths. It replaces the unexported loginHandler.
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) < 4 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth path not found")
		return
	}

	action := segs[2]
	provider := segs[3]
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return provider, nil
	}

	switch action {
	case "login":
		gothic.BeginAuthHandler(w, r)
	case "callback":
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cUser := &userDomain.ChatUser{User: user}
		cUser.SetUniqueID(getUserID(user.Email))
		avatarUrl, err := avatar.Avatars.GetAvatarURL(cUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := map[string]string{
			"user_id":    cUser.UniqueID(),
			"name":       user.Name,
			"email":      user.Email,
			"avatar_url": avatarUrl,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		authCookieValue := base64.StdEncoding.EncodeToString(jsonData)
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})

		if err := h.userService.Create(r.Context(), userDomain.User{
			ID:        cUser.UniqueID(),
			Name:      user.Name,
			Email:     user.Email,
			AvatarURL: avatarUrl,
			CreatedAt: time.Now(),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", "http://localhost:8080/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

// LogoutHandler handles logging out
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header().Set("Location", "http://localhost:8080/login")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
