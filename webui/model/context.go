package model

type ContextDocument struct {
	ID         string  `json:"id"`
	ProjectID  *string `json:"project_id"`
	TimeCodeID *string `json:"time_code_id"`
	Content    string  `json:"content"`
	CreatedAt  Time    `json:"created_at"`
	UpdatedAt  Time    `json:"updated_at"`
}
