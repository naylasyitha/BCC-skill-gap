package dto

type SelfAssessmentRequest struct {
	CareerID string              `json:"career_id" binding:"required,uuid"`
	Skills   []SkillLevelRequest `json:"skills" binding:"required,min=1,dive"`
}

type SkillLevelRequest struct {
	SkillID   string `json:"skill_id" binding:"required,uuid"`
	UserLevel string `json:"user_level" binding:"required"`
}

type SelfAssessmentResponse struct {
	UserCareerSessionID string `json:"user_career_session_id"`
}
