package storage

import "time"

type Entry struct {
	At    time.Time
	Level string
	Text  string
}
