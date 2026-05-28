package message

type Message struct {
	Text       string
	Hidden     []bool
	Converted  bool
	Persistent bool
	X          int
	Y          int
}

func New(text string, width, height int) *Message {
	return &Message{
		Text:   text,
		Hidden: make([]bool, len(text)),
		X:      (width - len(text)) / 2,
		Y:      height / 2,
	}
}

func (m *Message) SetText(text string, width, height int) {
	m.Text = text
	m.Hidden = make([]bool, len(text))
	m.Converted = false
	m.Persistent = false
	m.UpdatePosition(width, height)
}

func (m *Message) UpdatePosition(width, height int) {
	m.X = (width - len(m.Text)) / 2
	m.Y = height / 2
}
