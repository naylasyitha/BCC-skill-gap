package dto

type SkillCreateRequest struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc" binding:"required"`
}

type SkillResponse struct {
	ID   string `json:"skill_id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type SkillUpdateRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type CareerSkillCreateRequest struct {
	CareerID      string `json:"career_id" binding:"required"`
	SkillID       string `json:"skill_id" binding:"required"`
	Priority      int    `json:"priority" binding:"required"`
	RequiredLevel string `json:"required_level" binding:"required"`
}

type CareerSkillAsignResponse struct {
	ID            string `json:"id"`
	CareerID      string `json:"career_id"`
	SkillID       string `json:"skill_id"`
	Priority      int    `json:"priority"`
	RequiredLevel string `json:"required_level"`
}

type CareerSkillUpdateRequest struct {
	Priority      int    `json:"priority"`
	RequiredLevel string `json:"required_level"`
}
