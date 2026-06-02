package notification

import (
	"database/sql"
	"strings"
	"time"
)

func loadMetaValue(path, key string) (string, bool, error) {
	var value string
	found := false
	err := withLockedDB(path, func(db *sql.DB) error {
		err := db.QueryRow("SELECT value FROM notification_meta WHERE key = ?", strings.TrimSpace(key)).Scan(&value)
		if err == sql.ErrNoRows {
			return nil
		}
		if err != nil {
			return err
		}
		found = true
		return nil
	})
	if err != nil {
		return "", false, err
	}
	if !found {
		return "", false, nil
	}
	return strings.TrimSpace(value), true, nil
}

func saveMetaValue(path, key, value string) error {
	return withLockedDB(path, func(db *sql.DB) error {
		_, err := db.Exec(
			"INSERT INTO notification_meta (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value",
			strings.TrimSpace(key),
			strings.TrimSpace(value),
		)
		return err
	})
}

func loadMetaTime(path, key string) (time.Time, bool, error) {
	value, ok, err := loadMetaValue(path, key)
	if err != nil || !ok {
		return time.Time{}, ok, err
	}
	parsed, parseErr := time.Parse(time.RFC3339Nano, value)
	if parseErr != nil {
		return time.Time{}, false, nil
	}
	return parsed.UTC(), true, nil
}

func saveMetaTime(path, key string, at time.Time) error {
	return saveMetaValue(path, key, at.UTC().Format(time.RFC3339Nano))
}

func loadMetaBool(path, key string) (bool, bool, error) {
	value, ok, err := loadMetaValue(path, key)
	if err != nil || !ok {
		return false, ok, err
	}
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "true" {
		return true, true, nil
	}
	if value == "false" {
		return false, true, nil
	}
	return false, false, nil
}

func saveMetaBool(path, key string, value bool) error {
	if value {
		return saveMetaValue(path, key, "true")
	}
	return saveMetaValue(path, key, "false")
}
