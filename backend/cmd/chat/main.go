package main

import (
	"flag"
	"log"
	"os"

	"github.com/fhiroki/chat/internal/domain/message"
	"github.com/fhiroki/chat/internal/domain/user"
	"github.com/fhiroki/chat/internal/infrastructure/db"
	"github.com/fhiroki/chat/internal/infrastructure/db/sqlite"
	"github.com/fhiroki/chat/internal/infrastructure/websocket"
	"github.com/fhiroki/chat/internal/interfaces/handler"
	"github.com/fhiroki/chat/internal/interfaces/middleware"
	"github.com/fhiroki/chat/internal/trace"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

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

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/avatars", "./avatars")

	// ミドルウェアの設定
	auth := r.Group("/", middleware.MustAuth())

	// ルーティングの設定
	auth.GET("/chat", handler.TemplateHandler("chat.html"))
	auth.GET("/upload", handler.TemplateHandler("upload.html"))
	auth.POST("/uploader", handler.UploaderHandler)
	auth.GET("/messages", messageHandler.GetMessages)
	auth.POST("/messages", messageHandler.CreateMessage)
	auth.DELETE("/messages/:id", messageHandler.DeleteMessage)
	auth.GET("/room", room.HandleWebSocket)

	r.GET("/login", handler.TemplateHandler("login.html"))
	r.GET("/logout", authHandler.LogoutHandler)
	r.GET("/auth/*provider", authHandler.LoginHandler)

	go room.Run()
	if err := r.Run(*addr); err != nil {
		log.Fatal(err)
	}
}
