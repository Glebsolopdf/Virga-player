package message

import "unicode/utf8"

type Message struct {
	Text       string
	Hidden     []bool
	Converted  bool
	Persistent bool
	X          int
	Y          int
}

func New(text string, width, height int) *Message {
	runeCount := utf8.RuneCountInString(text)
	return &Message{
		Text:   text,
		Hidden: make([]bool, runeCount),
		X:      (width - runeCount) / 2,
		Y:      height / 2,
	}
}

func (m *Message) SetText(text string, width, height int) {
	m.Text = text
	m.Hidden = make([]bool, utf8.RuneCountInString(text))
	m.Converted = false
	m.Persistent = false
	m.UpdatePosition(width, height)
}

func (m *Message) UpdatePosition(width, height int) {
	m.X = (width - utf8.RuneCountInString(m.Text)) / 2
	m.Y = height / 2
}
