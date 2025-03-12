package goasteroids

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Timer is the type for all in-game timers. It holds the current tick count, and the target tick count.
type Timer struct {
	currentTicks int
	targetTicks  int
}

// NewTimer is a factory method for creating timers of duration d.
func NewTimer(d time.Duration) *Timer {
	return &Timer {
		currentTicks: 0,
		targetTicks: int(d.Milliseconds()) * ebiten.TPS() / 1000,
	}
}

// Update updates the timer's currentTicks count.
func (t *Timer) Update() {
	if t.currentTicks < t.targetTicks {
		t.currentTicks++
	}
}

// IsReady returns true if the timer's target has been reach, and otherwise false.
func (t *Timer) IsReady() bool {
	return t.currentTicks >= t.targetTicks
}

// Reset sets a timer back to zero.
func (t *Timer) Reset() {
	t.currentTicks = 0
}
