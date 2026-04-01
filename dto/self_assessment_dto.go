package dto

type SelfAssessmentRequest struct {
	Skills []SkillLevelRequest `json:"skills" binding:"required,min=1,dive"`
}

type SkillLevelRequest struct {
	SkillID   string `json:"skill_id" binding:"required,uuid"`
	UserLevel string `json:"user_level"`
}

type SelfAssessmentResponse struct {
	UserCareerSessionID string `json:"career_session_id"`
}
