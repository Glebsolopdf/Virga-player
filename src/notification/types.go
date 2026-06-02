package notification

import "time"

type Kind string

const (
	KindInfo    Kind = "info"
	KindWelcome Kind = "welcome"
	KindUpdate  Kind = "update"
)

type Item struct {
	ID        string     `json:"id"`
	Kind      Kind       `json:"kind"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Version   string     `json:"version,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

func (i Item) IsRead() bool {
	return i.ReadAt != nil
}

type State struct {
	Items []Item `json:"items"`
}
