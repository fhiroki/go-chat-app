package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Driver string
	DSN    string
}

func NewConnection(config Config) (*sql.DB, error) {
	db, err := sql.Open(config.Driver, config.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func MustNewConnection(config Config) *sql.DB {
	db, err := NewConnection(config)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return db
}

// InitializeTables はデータベースの必要なテーブルを作成します
func InitializeTables(db *sql.DB) error {
	// SQLファイルからマイグレーションを実行
	migrationPath := filepath.Join("internal", "infrastructure", "db", "migrations", "001_create_tables.sql")
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// SQLファイルを実行
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	log.Println("Database tables initialized")
	return nil
}
