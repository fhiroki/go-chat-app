package handler

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

type TemplateHandler struct {
	once     sync.Once
	Filename string
	templ    *template.Template
}

func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.Filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		decoded, err := base64.StdEncoding.DecodeString(authCookie.Value)
		if err == nil {
			var userData map[string]string
			if jsonErr := json.Unmarshal(decoded, &userData); jsonErr == nil {
				data["UserData"] = userData
			}
		}
	}
	t.templ.Execute(w, data)
}
