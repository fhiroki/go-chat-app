package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/markbates/goth/gothic"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler は /auth/{action}/{provider} の URL パスにマッチする認証処理ハンドラです。
// リクエストの action に "login" が指定された場合、認証を開始し、
// "callback" の場合は認証のコールバックを処理してクライアントに認証情報をセットします。
func loginHandler(w http.ResponseWriter, r *http.Request) {
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

		// なぜかNameが取得できない
		data := map[string]string{
			"name":       user.Name,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
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

		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
