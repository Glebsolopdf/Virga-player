package storage

import "strings"

func (b *RingBuffer) Dump() string {
	items := b.Last(b.cap)
	if len(items) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, e := range items {
		sb.WriteString(e.At.Format("2026-01-02 15:04:05"))
		sb.WriteString(" [")
		sb.WriteString(e.Level)
		sb.WriteString("] ")
		sb.WriteString(e.Text)
		sb.WriteByte('\n')
	}
	return sb.String()
}
