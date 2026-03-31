package dto

type CareerSessionCreateRequest struct {
	CareerID string `json:"career_id" binding:"required,uuid"`
}

type CareerSessionResponse struct {
	ID          string `json:"career_session_id"`
	UserID      string `json:"user_id"`
	CareerID    string `json:"career_id"`
	Status      string `json:"status"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

type CareerSessionDetailResponse struct {
	ID          string `json:"career_session_id"`
	UserID      string `json:"user_id"`
	Fullname    string `json:"full_name"`
	CareerID    string `json:"career_id"`
	CareerName  string `json:"career_name"`
	Status      string `json:"status"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}
