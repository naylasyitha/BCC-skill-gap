package dto

type CareerAnalyticResponse struct {
	CareerSessionID string                  `json:"career_session_id"`
	TotalScore      int                     `json:"total_score"`
	SkillsResult    []SkillAnalyticResponse `json:"skills_result"`
	Recommendations []SkillRecommendation   `json:"recommendations"`
}
type SkillRecommendation struct {
	SkillID      string `json:"skill_id"`
	SkillName    string `json:"skill_name"`
	CurrentLevel string `json:"current_level"`
	TargetLevel  string `json:"target_level"`
	Priority     int    `json:"priority"`
}

type SkillAnalyticResponse struct {
	SkillID         string   `json:"skill_id"`
	SkillName       string   `json:"skill_name"`
	UserLevel       string   `json:"user_level"`
	FinalLevel      string   `json:"final_user_level"`
	RequiredLevel   string   `json:"required_level"`
	SkillScore      int      `json:"skill_score"`
	GapLevel        int      `json:"gap_level"`
	Status          string   `json:"status"`
	SuggestionLevel []string `json:"suggestion_level"`
}
