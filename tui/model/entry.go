package model

type Entry struct {
	ID                 string   `json:"id"`
	RawText            string   `json:"raw_text"`
	DurationMinutes    *int     `json:"duration_minutes"`
	WorkTypeID         *string  `json:"work_type_id"`
	TimeCodeID         *string  `json:"time_code_id"`
	Description        *string  `json:"description"`
	EntryDate          string   `json:"entry_date"`
	ConfidenceDuration *float64 `json:"confidence_duration"`
	ConfidenceWorkType *float64 `json:"confidence_work_type"`
	ConfidenceTimeCode *float64 `json:"confidence_time_code"`
	ConfidenceOverall  *float64 `json:"confidence_overall"`
	Status             string   `json:"status"`
	ParseError         *string  `json:"parse_error"`
	ManuallyCorrect    bool     `json:"manually_corrected"`
	Locked             bool     `json:"locked"`
	CorrectedAt        *Time    `json:"corrected_at"`
	CreatedAt          Time     `json:"created_at"`
	ParsedAt           *Time    `json:"parsed_at"`
	UpdatedAt          Time     `json:"updated_at"`
}
