package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/fhiroki/chat/internal/domain/message"
	"github.com/fhiroki/chat/internal/domain/user"
	"github.com/fhiroki/chat/internal/infrastructure/db"
	"github.com/fhiroki/chat/internal/infrastructure/db/sqlite"
	"github.com/fhiroki/chat/internal/infrastructure/websocket"
	"github.com/fhiroki/chat/internal/interfaces/handler"
	"github.com/fhiroki/chat/internal/interfaces/middleware"
	"github.com/fhiroki/chat/internal/trace"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
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

func main() {
	addr := flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	// データベース接続の設定
	dbConn := db.MustNewConnection(db.Config{
		Driver: "sqlite3",
		DSN:    "chat.db",
	})
	defer dbConn.Close()

	// テーブルの初期化
	if err := db.InitializeTables(dbConn); err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	txManager := db.NewTxManager(dbConn)

	// メッセージ関連の初期化
	messageRepo := sqlite.NewMessageRepository(dbConn)
	messageService := message.NewMessageService(messageRepo, txManager)
	messageHandler := handler.NewMessageHandler(messageService)

	// ユーザー関連の初期化
	userRepo := sqlite.NewUserRepository(dbConn)
	userService := user.NewUserService(userRepo, txManager)
	authHandler := handler.NewAuthHandler(userService)

	// WebSocket関連の初期化
	room := websocket.NewRoom(messageService)
	room.Tracer = trace.New(os.Stdout)

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), os.Getenv("GOOGLE_REDIRECT_URL")),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), os.Getenv("GITHUB_REDIRECT_URL")),
	)

	http.Handle("/chat", middleware.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.HandleFunc("/uploader", handler.UploaderHandler)
	http.HandleFunc("/logout", authHandler.LogoutHandler)
	http.HandleFunc("/auth/", authHandler.LoginHandler)
	http.HandleFunc("/messages", messageHandler.HandleMessages)
	http.HandleFunc("/messages/delete", messageHandler.DeleteMessage)
	http.Handle("/room", room)

	go room.Run()
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
