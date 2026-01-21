package timer

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manager handles timer lifecycle and persistence
type Manager struct {
	timers      []*Timer
	persistPath string
}

// NewManager creates a new timer manager and loads persisted timers
func NewManager() *Manager {
	m := &Manager{
		timers:      make([]*Timer, 0),
		persistPath: defaultPersistPath(),
	}
	m.load()
	m.cleanupOldCompleted()
	return m
}

// defaultPersistPath returns the path to the timers JSON file
func defaultPersistPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "timers.json"
	}
	return filepath.Join(home, ".tact", "timers.json")
}

// StartTimer creates and starts a new timer, pausing any running timer
func (m *Manager) StartTimer(description string) *Timer {
	// Pause any running timer first
	m.pauseRunning()

	t := NewTimer(generateID(), description)
	m.timers = append(m.timers, t)
	m.save()
	return t
}

// generateID creates a simple unique ID
func generateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// PauseTimer pauses the timer with the given ID
func (m *Manager) PauseTimer(id string) {
	for _, t := range m.timers {
		if t.ID == id {
			t.Pause()
			m.save()
			return
		}
	}
}

// ResumeTimer resumes the timer with the given ID, pausing any other running timer
func (m *Manager) ResumeTimer(id string) {
	// Pause any running timer first
	m.pauseRunning()

	for _, t := range m.timers {
		if t.ID == id {
			t.Resume()
			m.save()
			return
		}
	}
}

// StopTimer stops the timer with the given ID and returns it
func (m *Manager) StopTimer(id string) *Timer {
	for _, t := range m.timers {
		if t.ID == id {
			t.Stop()
			m.save()
			return t
		}
	}
	return nil
}

// DeleteTimer removes the timer with the given ID
func (m *Manager) DeleteTimer(id string) {
	for i, t := range m.timers {
		if t.ID == id {
			m.timers = append(m.timers[:i], m.timers[i+1:]...)
			m.save()
			return
		}
	}
}

// RunningTimer returns the currently running timer, if any
func (m *Manager) RunningTimer() *Timer {
	for _, t := range m.timers {
		if t.IsRunning() {
			return t
		}
	}
	return nil
}

// ActiveTimers returns all running or paused timers
func (m *Manager) ActiveTimers() []*Timer {
	var active []*Timer
	for _, t := range m.timers {
		if t.IsRunning() || t.IsPaused() {
			active = append(active, t)
		}
	}
	return active
}

// localMidnight returns the start of today in local timezone
func localMidnight() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

// CompletedToday returns all stopped timers from today
func (m *Manager) CompletedToday() []*Timer {
	var completed []*Timer
	today := localMidnight()

	for _, t := range m.timers {
		if t.IsStopped() && t.StoppedAt != nil {
			if !t.StoppedAt.Before(today) {
				completed = append(completed, t)
			}
		}
	}
	return completed
}

// GetTimer returns the timer with the given ID
func (m *Manager) GetTimer(id string) *Timer {
	for _, t := range m.timers {
		if t.ID == id {
			return t
		}
	}
	return nil
}

// pauseRunning pauses any currently running timer
func (m *Manager) pauseRunning() {
	for _, t := range m.timers {
		if t.IsRunning() {
			t.Pause()
		}
	}
}

// cleanupOldCompleted removes completed timers from previous days
func (m *Manager) cleanupOldCompleted() {
	today := localMidnight()
	var kept []*Timer

	for _, t := range m.timers {
		// Keep all active timers
		if !t.IsStopped() {
			kept = append(kept, t)
			continue
		}
		// Keep completed timers from today (stopped at or after local midnight)
		if t.StoppedAt != nil && !t.StoppedAt.Before(today) {
			kept = append(kept, t)
		}
	}

	if len(kept) != len(m.timers) {
		m.timers = kept
		m.save()
	}
}

// load reads timers from the persistence file
func (m *Manager) load() {
	data, err := os.ReadFile(m.persistPath)
	if err != nil {
		return // File doesn't exist yet, start fresh
	}

	var timers []*Timer
	if err := json.Unmarshal(data, &timers); err != nil {
		return // Invalid JSON, start fresh
	}
	m.timers = timers
}

// save writes timers to the persistence file
func (m *Manager) save() {
	// Ensure directory exists
	dir := filepath.Dir(m.persistPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}

	data, err := json.MarshalIndent(m.timers, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(m.persistPath, data, 0644)
}
