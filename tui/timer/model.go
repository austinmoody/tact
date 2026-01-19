package timer

import (
	"time"
)

// TimerState represents the current state of a timer
type TimerState string

const (
	TimerRunning TimerState = "running"
	TimerPaused  TimerState = "paused"
	TimerStopped TimerState = "stopped"
)

// Timer represents a time tracking timer
type Timer struct {
	ID                 string     `json:"id"`
	Description        string     `json:"description"`
	State              TimerState `json:"state"`
	StartedAt          *time.Time `json:"started_at,omitempty"`
	AccumulatedSeconds int        `json:"accumulated_seconds"`
	StoppedAt          *time.Time `json:"stopped_at,omitempty"`
}

// NewTimer creates a new running timer with the given description
func NewTimer(id, description string) *Timer {
	now := time.Now()
	return &Timer{
		ID:                 id,
		Description:        description,
		State:              TimerRunning,
		StartedAt:          &now,
		AccumulatedSeconds: 0,
	}
}

// TotalElapsedSeconds returns the total elapsed time including current running period
func (t *Timer) TotalElapsedSeconds() int {
	total := t.AccumulatedSeconds
	if t.State == TimerRunning && t.StartedAt != nil {
		total += int(time.Since(*t.StartedAt).Seconds())
	}
	return total
}

// Pause pauses a running timer, accumulating the elapsed time
func (t *Timer) Pause() {
	if t.State != TimerRunning || t.StartedAt == nil {
		return
	}
	t.AccumulatedSeconds += int(time.Since(*t.StartedAt).Seconds())
	t.StartedAt = nil
	t.State = TimerPaused
}

// Resume resumes a paused timer
func (t *Timer) Resume() {
	if t.State != TimerPaused {
		return
	}
	now := time.Now()
	t.StartedAt = &now
	t.State = TimerRunning
}

// Stop stops a timer, accumulating any remaining running time
func (t *Timer) Stop() {
	if t.State == TimerRunning && t.StartedAt != nil {
		t.AccumulatedSeconds += int(time.Since(*t.StartedAt).Seconds())
	}
	t.StartedAt = nil
	t.State = TimerStopped
	now := time.Now()
	t.StoppedAt = &now
}

// IsRunning returns true if the timer is currently running
func (t *Timer) IsRunning() bool {
	return t.State == TimerRunning
}

// IsPaused returns true if the timer is paused
func (t *Timer) IsPaused() bool {
	return t.State == TimerPaused
}

// IsStopped returns true if the timer is stopped
func (t *Timer) IsStopped() bool {
	return t.State == TimerStopped
}
