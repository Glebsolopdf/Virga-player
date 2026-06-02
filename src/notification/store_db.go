package notification

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

var notificationDBMu sync.Mutex

func openDB(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA busy_timeout = 30000"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA synchronous = NORMAL"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := initSchema(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func withLockedDB(path string, fn func(*sql.DB) error) error {
	notificationDBMu.Lock()
	defer notificationDBMu.Unlock()

	db, err := openDB(path)
	if err != nil {
		return err
	}
	defer db.Close()

	return fn(db)
}

func initSchema(db *sql.DB) error {
	const schema = `
CREATE TABLE IF NOT EXISTS notifications (
  id TEXT PRIMARY KEY,
  kind TEXT NOT NULL,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  version TEXT,
  created_at TEXT NOT NULL,
  read_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

CREATE TABLE IF NOT EXISTS notification_meta (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL
);
`
	_, err := db.Exec(schema)
	return err
}

func pruneExpiredRows(db *sql.DB, now time.Time) error {
	cutoff := now.UTC().Add(-notificationMaxAge).Format(time.RFC3339Nano)
	_, err := db.Exec("DELETE FROM notifications WHERE created_at < ?", cutoff)
	return err
}

func pruneExpiredRowsTx(tx *sql.Tx, now time.Time) error {
	cutoff := now.UTC().Add(-notificationMaxAge).Format(time.RFC3339Nano)
	_, err := tx.Exec("DELETE FROM notifications WHERE created_at < ?", cutoff)
	return err
}
