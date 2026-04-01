package dto

type CareerAnalyticResponse struct {
	CareerSessionID string                  `json:"career_session_id"`
	TotalScore      int                     `json:"total_score"`
	SkillsResult    []SkillAnalyticResponse `json:"skills_result"`
}

type SkillAnalyticResponse struct {
	SkillID    string `json:"skill_id"`
	SkillName  string `json:"skill_name"`
	UserLevel  string `json:"user_level"`
	FinalLevel string `json:"final_user_level"`
	SkillScore int    `json:"skill_score"`
}
