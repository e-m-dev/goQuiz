package store

import (
	"context"
	"database/sql"
	"goQuiz/server/internal/cfg"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		if cfg.Debug {
			log.Printf("SQLITE -> OPEN | Error accessing db path")
		}
		return nil, err
	}

	d, err := sql.Open("sqlite", path)
	if err != nil {
		if cfg.Debug {
			log.Printf("SQLITE -> OPEN | Error opening db")
		}
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err = d.ExecContext(ctx, `
		PRAGMA journal_mode=WAL;
		PRAGMA synchronous=NORMAL;
		PRAGMA foreign_keys=ON;
		PRAGMA busy_timeout=5000;
	`); err != nil {
		if cfg.Debug {
			log.Printf("SQLITE -> OPEN | Error setting pragmas")
		}
		_ = d.Close()
		return nil, err
	}

	if err = d.PingContext(ctx); err != nil {
		_ = d.Close()
		return nil, err
	}

	if cfg.Debug {
		log.Printf("SQLITE -> OPEN | path = {%s}", path)
	}

	return d, nil
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS questions (
			id				INTEGER PRIMARY KEY AUTOINCREMENT,
			prompt			TEXT	NOT NULL,
			options_json	TEXT	NOT NULL,
			correct_index	INTEGER	NOT NULL CHECK (correct_index >= 0),
			category		TEXT,
			difficulty		TEXT,
			created_at		TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
			CREATE INDEX IF NOT EXISTS idx_questions_cat	ON questions(category);
			CREATE INDEX IF NOT EXISTS idx_questions_diff	ON questions(difficulty);
	`)
	if err == nil && cfg.Debug {
		log.Printf("SQLITE -> MIGRATE | ok")
	}
	return err
}
