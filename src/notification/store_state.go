package notification

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func loadState(path string) (State, error) {
	state := State{Items: []Item{}}
	err := withLockedDB(path, func(db *sql.DB) error {
		if err := migrateLegacyJSON(db, legacyStatePathFor(path)); err != nil {
			return err
		}
		if err := pruneExpiredRows(db, time.Now().UTC()); err != nil {
			return err
		}

		rows, err := db.Query(`
			SELECT id, kind, title, body, COALESCE(version, ''), created_at, read_at
			FROM notifications
			ORDER BY created_at DESC, id DESC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				id        string
				kind      string
				title     string
				body      string
				version   string
				createdAt string
				readAt    sql.NullString
			)
			if err := rows.Scan(&id, &kind, &title, &body, &version, &createdAt, &readAt); err != nil {
				return err
			}

			item := Item{
				ID:      strings.TrimSpace(id),
				Kind:    Kind(strings.TrimSpace(kind)),
				Title:   fallbackString(title, unsupportedTitle),
				Body:    fallbackString(body, unsupportedBodyText),
				Version: strings.TrimSpace(version),
			}

			parsedCreated, err := time.Parse(time.RFC3339Nano, createdAt)
			if err != nil {
				item.CreatedAt = time.Now().UTC()
				item.Title = unsupportedTitle
				item.Body = unsupportedBodyText
			} else {
				item.CreatedAt = parsedCreated
			}

			if readAt.Valid && strings.TrimSpace(readAt.String) != "" {
				if parsedRead, err := time.Parse(time.RFC3339Nano, readAt.String); err == nil {
					stamp := parsedRead
					item.ReadAt = &stamp
				}
			}

			if item.ID == "" {
				item.ID = fmt.Sprintf("legacy:%d", item.CreatedAt.UnixNano())
			}
			if item.Kind == "" {
				item.Kind = KindInfo
			}

			state.Items = append(state.Items, item)
		}
		return rows.Err()
	})
	if err != nil {
		return State{}, err
	}
	return state, nil
}

func saveState(path string, state State) error {
	return withLockedDB(path, func(db *sql.DB) error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer func() {
			_ = tx.Rollback()
		}()

		if _, err := tx.Exec("DELETE FROM notifications"); err != nil {
			return err
		}

		stmt, err := tx.Prepare(`
			INSERT INTO notifications (id, kind, title, body, version, created_at, read_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, item := range state.Items {
			id := strings.TrimSpace(item.ID)
			if id == "" {
				id = fmt.Sprintf("legacy:%d", item.CreatedAt.UnixNano())
			}
			kind := strings.TrimSpace(string(item.Kind))
			if kind == "" {
				kind = string(KindInfo)
			}
			createdAt := item.CreatedAt.UTC()
			if createdAt.IsZero() {
				createdAt = time.Now().UTC()
			}

			var readAt any
			if item.ReadAt != nil && !item.ReadAt.IsZero() {
				readAt = item.ReadAt.UTC().Format(time.RFC3339Nano)
			}

			if _, err := stmt.Exec(
				id,
				kind,
				fallbackString(item.Title, unsupportedTitle),
				fallbackString(item.Body, unsupportedBodyText),
				strings.TrimSpace(item.Version),
				createdAt.Format(time.RFC3339Nano),
				readAt,
			); err != nil {
				return err
			}
		}

		if err := pruneExpiredRowsTx(tx, time.Now().UTC()); err != nil {
			return err
		}

		return tx.Commit()
	})
}

func clearState(path string) error {
	state := State{Items: []Item{}}
	return saveState(path, state)
}
