package animation

import (
	"time"
)

// Engine manages animation timing and FPS
type Engine struct {
	ticker    *time.Ticker
	frameTime time.Duration
}

// NewEngine creates a new animation engine with target FPS
func NewEngine(targetFPS int) *Engine {
	frameTime := time.Second / time.Duration(targetFPS)
	return &Engine{
		ticker:    time.NewTicker(frameTime),
		frameTime: frameTime,
	}
}

// Tick returns a channel that sends on each frame
func (e *Engine) Tick() <-chan time.Time {
	return e.ticker.C
}

// FrameDuration returns the target duration between frames.
func (e *Engine) FrameDuration() time.Duration {
	return e.frameTime
}

// Stop stops the animation engine
func (e *Engine) Stop() {
	e.ticker.Stop()
}
