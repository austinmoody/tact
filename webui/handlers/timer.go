package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	"tact-webui/api"
	"tact-webui/templates/pages"
)

type TimerHandler struct {
	client *api.Client
	state  TimerStateData
	mu     sync.RWMutex
}

type TimerStateData struct {
	Active      bool
	Paused      bool
	Description string
	StartTime   time.Time
	PausedAt    time.Time
	PausedTotal time.Duration
}

func (h *TimerHandler) getState() pages.TimerState {
	h.mu.RLock()
	defer h.mu.RUnlock()

	elapsed := "00:00:00"
	if h.state.Active {
		var d time.Duration
		if h.state.Paused {
			d = h.state.PausedAt.Sub(h.state.StartTime) - h.state.PausedTotal
		} else {
			d = time.Since(h.state.StartTime) - h.state.PausedTotal
		}
		hours := int(d.Hours())
		mins := int(d.Minutes()) % 60
		secs := int(d.Seconds()) % 60
		elapsed = fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)
	}

	return pages.TimerState{
		Active:      h.state.Active,
		Paused:      h.state.Paused,
		Description: h.state.Description,
		Elapsed:     elapsed,
	}
}

func (h *TimerHandler) Page(w http.ResponseWriter, r *http.Request) {
	pages.Timer(h.getState()).Render(r.Context(), w)
}

func (h *TimerHandler) Stream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			state := h.getState()

			// Send timer display update
			fmt.Fprintf(w, "event: timer-display\ndata: %s\n\n", state.Elapsed)

			// Send timer indicator for nav bar
			var buf bytes.Buffer
			pages.TimerIndicator(state.Active, state.Elapsed).Render(r.Context(), &buf)
			fmt.Fprintf(w, "event: timer-tick\ndata: %s\n\n", buf.String())

			flusher.Flush()
		}
	}
}

func (h *TimerHandler) Start(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")
	if description == "" {
		description = "Untitled timer"
	}

	h.mu.Lock()
	h.state = TimerStateData{
		Active:      true,
		Paused:      false,
		Description: description,
		StartTime:   time.Now(),
	}
	h.mu.Unlock()

	// For HTMX requests, use HX-Redirect for full page navigation
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/timer")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/timer", http.StatusSeeOther)
}

func (h *TimerHandler) Pause(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	if h.state.Active && !h.state.Paused {
		h.state.Paused = true
		h.state.PausedAt = time.Now()
	}
	h.mu.Unlock()

	pages.TimerControls(h.getState()).Render(r.Context(), w)
}

func (h *TimerHandler) Resume(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	if h.state.Active && h.state.Paused {
		h.state.PausedTotal += time.Since(h.state.PausedAt)
		h.state.Paused = false
	}
	h.mu.Unlock()

	pages.TimerControls(h.getState()).Render(r.Context(), w)
}

func (h *TimerHandler) Stop(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	description := h.state.Description
	var elapsed time.Duration
	if h.state.Active {
		if h.state.Paused {
			elapsed = h.state.PausedAt.Sub(h.state.StartTime) - h.state.PausedTotal
		} else {
			elapsed = time.Since(h.state.StartTime) - h.state.PausedTotal
		}
	}
	h.state = TimerStateData{}
	h.mu.Unlock()

	// Create entry from timer
	if elapsed > 0 {
		minutes := int(elapsed.Minutes())
		if minutes < 1 {
			minutes = 1 // Minimum 1 minute
		}
		userInput := fmt.Sprintf("%dm %s", minutes, description)
		h.client.CreateEntry(userInput)
	}

	// For HTMX requests, use HX-Redirect for full page navigation
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/timer")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/timer", http.StatusSeeOther)
}

func (h *TimerHandler) Discard(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.state = TimerStateData{}
	h.mu.Unlock()

	// For HTMX requests, use HX-Redirect for full page navigation
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/timer")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/timer", http.StatusSeeOther)
}
