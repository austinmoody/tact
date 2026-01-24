package model

import "fmt"

type Entry struct {
	ID                 string   `json:"id"`
	UserInput          string   `json:"user_input"`
	DurationMinutes    *int     `json:"duration_minutes"`
	WorkTypeID         *string  `json:"work_type_id"`
	TimeCodeID         *string  `json:"time_code_id"`
	ParsedDescription  *string  `json:"parsed_description"`
	EntryDate          string   `json:"entry_date"`
	ConfidenceDuration *float64 `json:"confidence_duration"`
	ConfidenceWorkType *float64 `json:"confidence_work_type"`
	ConfidenceTimeCode *float64 `json:"confidence_time_code"`
	ConfidenceOverall  *float64 `json:"confidence_overall"`
	Status             string   `json:"status"`
	ParseError         *string  `json:"parse_error"`
	ParseNotes         *string  `json:"parse_notes"`
	ManuallyCorrect    bool     `json:"manually_corrected"`
	Locked             bool     `json:"locked"`
	CorrectedAt        *Time    `json:"corrected_at"`
	CreatedAt          Time     `json:"created_at"`
	ParsedAt           *Time    `json:"parsed_at"`
	UpdatedAt          Time     `json:"updated_at"`
}

// StatusColor returns the CSS class for the entry status
func (e Entry) StatusColor() string {
	switch e.Status {
	case "parsed":
		return "success"
	case "pending":
		return "warning"
	case "failed":
		return "error"
	default:
		return ""
	}
}

// DurationDisplay returns a formatted duration string
func (e Entry) DurationDisplay() string {
	if e.DurationMinutes == nil {
		return "-"
	}
	hours := *e.DurationMinutes / 60
	mins := *e.DurationMinutes % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}
