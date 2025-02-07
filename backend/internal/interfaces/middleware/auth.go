package middleware

import (
	"net/http"
)

// authHandler handles authentication middleware
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

// MustAuth ensures authentication
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
