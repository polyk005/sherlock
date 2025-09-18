package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteConfig struct {
	Path string
}

func NewSQLiteDB(cfg SQLiteConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия БД: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	// Прагма для лучшей производительности
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA synchronous=NORMAL;")
	db.Exec("PRAGMA foreign_keys=ON;")

	return db, nil
}

func CreateSubscribersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS subscribers (
		user_id INTEGER PRIMARY KEY,
		username TEXT,
		first_name TEXT,
		last_name TEXT,
		joined_at DATETIME,
		last_check DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_subscribers_user_id ON subscribers(user_id);
	CREATE INDEX IF NOT EXISTS idx_subscribers_joined_at ON subscribers(joined_at);
	`

	_, err := db.Exec(query)
	return err
}
