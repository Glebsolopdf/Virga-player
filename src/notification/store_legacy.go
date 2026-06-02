package notification

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

func migrateLegacyJSON(db *sql.DB, legacyPath string) error {
	data, err := os.ReadFile(legacyPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	var raw struct {
		Items []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		// Leave the file in place to avoid destructive behavior on malformed data.
		return nil
	}

	if len(raw.Items) == 0 {
		_ = os.Remove(legacyPath)
		return nil
	}

	for index, payload := range raw.Items {
		item := decodeLegacyItem(payload, index)
		if _, err := db.Exec(
			`INSERT OR IGNORE INTO notifications (id, kind, title, body, version, created_at, read_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			item.ID,
			string(item.Kind),
			item.Title,
			item.Body,
			item.Version,
			item.CreatedAt.UTC().Format(time.RFC3339Nano),
			nullableReadAt(item.ReadAt),
		); err != nil {
			return err
		}
	}

	_ = os.Remove(legacyPath)
	return nil
}

func decodeLegacyItem(payload json.RawMessage, index int) Item {
	var mapValue map[string]any
	if err := json.Unmarshal(payload, &mapValue); err != nil {
		now := time.Now().UTC()
		return Item{
			ID:        fmt.Sprintf("legacy-unsupported:%d", index),
			Kind:      KindInfo,
			Title:     unsupportedTitle,
			Body:      unsupportedBodyText,
			CreatedAt: now,
		}
	}

	item := Item{
		ID:      stringField(mapValue, "id", fmt.Sprintf("legacy:%d", index)),
		Kind:    Kind(stringField(mapValue, "kind", string(KindInfo))),
		Title:   stringField(mapValue, "title", unsupportedTitle),
		Body:    stringField(mapValue, "body", unsupportedBodyText),
		Version: stringField(mapValue, "version", ""),
	}

	if created, ok := timeField(mapValue, "created_at"); ok {
		item.CreatedAt = created
	} else {
		item.CreatedAt = time.Now().UTC()
		item.Title = unsupportedTitle
		item.Body = unsupportedBodyText
	}

	if readAt, ok := timeField(mapValue, "read_at"); ok {
		stamp := readAt
		item.ReadAt = &stamp
	}

	if strings.TrimSpace(item.ID) == "" {
		item.ID = fmt.Sprintf("legacy:%d", item.CreatedAt.UnixNano())
	}
	if item.Kind == "" {
		item.Kind = KindInfo
	}

	return item
}

func stringField(values map[string]any, key, fallback string) string {
	v, ok := values[key]
	if !ok {
		return fallback
	}
	text, ok := v.(string)
	if !ok {
		return fallback
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return fallback
	}
	return text
}

func timeField(values map[string]any, key string) (time.Time, bool) {
	v, ok := values[key]
	if !ok {
		return time.Time{}, false
	}
	text, ok := v.(string)
	if !ok {
		return time.Time{}, false
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return time.Time{}, false
	}
	parsed, err := time.Parse(time.RFC3339Nano, text)
	if err != nil {
		return time.Time{}, false
	}
	return parsed.UTC(), true
}

func nullableReadAt(readAt *time.Time) any {
	if readAt == nil || readAt.IsZero() {
		return nil
	}
	return readAt.UTC().Format(time.RFC3339Nano)
}

func fallbackString(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
