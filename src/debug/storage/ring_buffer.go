package storage

import "sync"

type RingBuffer struct {
	mu      sync.RWMutex
	entries []Entry
	next    int
	count   int
	cap     int
}

func NewRingBuffer(capacity int) *RingBuffer {
	if capacity < 20 {
		capacity = 20
	}
	return &RingBuffer{entries: make([]Entry, capacity), cap: capacity}
}

func (b *RingBuffer) Append(e Entry) {
	b.mu.Lock()
	b.entries[b.next] = e
	b.next = (b.next + 1) % b.cap
	if b.count < b.cap {
		b.count++
	}
	b.mu.Unlock()
}

func (b *RingBuffer) Last(n int) []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if n <= 0 || b.count == 0 {
		return nil
	}
	if n > b.count {
		n = b.count
	}
	out := make([]Entry, 0, n)
	start := b.next - n
	if start < 0 {
		start += b.cap
	}
	for i := 0; i < n; i++ {
		out = append(out, b.entries[(start+i)%b.cap])
	}
	return out
}
