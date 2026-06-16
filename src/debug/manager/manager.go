package manager

import (
	"sync"
	"time"

	"virga-player/debug/storage"
)

type Manager struct {
	mu           sync.RWMutex
	enabled      bool
	forced       bool
	buf          *storage.RingBuffer
	status       string
	statusAt     time.Time
	copyBtn      rect
	saveBtn      rect
	lastOverlayW int

	fps           float64
	particles     int
	particlesMax  int
	targetFPS     int
	cpuPercent    float64
	memMiB        float64
	goroutines    int
	lastCPUTime   float64
	lastCPUSample time.Time
}
