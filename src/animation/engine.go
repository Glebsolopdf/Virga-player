package animation

import (
	"time"
)

// Engine manages animation timing and FPS
type Engine struct {
	ticker     *time.Ticker
	fpsTarget  int
	frameTime  time.Duration
	lastFPS    float64
	frameCount int
	lastTime   time.Time
}

// NewEngine creates a new animation engine with target FPS
func NewEngine(targetFPS int) *Engine {
	frameTime := time.Second / time.Duration(targetFPS)
	return &Engine{
		ticker:    time.NewTicker(frameTime),
		fpsTarget: targetFPS,
		frameTime: frameTime,
		lastTime:  time.Now(),
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

// FPS returns the current FPS
func (e *Engine) FPS() float64 {
	now := time.Now()
	elapsed := now.Sub(e.lastTime).Seconds()

	if elapsed >= 1.0 {
		e.lastFPS = float64(e.frameCount) / elapsed
		e.frameCount = 0
		e.lastTime = now
	}

	e.frameCount++
	return e.lastFPS
}

// Stop stops the animation engine
func (e *Engine) Stop() {
	e.ticker.Stop()
}
