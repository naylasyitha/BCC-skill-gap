package dto

type CareerCreateRequest struct {
	Name   string               `json:"name" binding:"required"`
	Desc   string               `json:"desc" binding:"required"`
	Skills []CareerSkillRequest `json:"skills" binding:"required,min=1,dive"`
}

type CareerSkillRequest struct {
	SkillID       string `json:"skill_id" binding:"required,uuid"`
	Priority      int    `json:"priority" binding:"required"`
	RequiredLevel string `json:"required_level" binding:"required"`
}

type CareerResponse struct {
	ID   string `json:"career_id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type SkillsResponse struct {
	ID            string `json:"skill_id"`
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Priority      int    `json:"priority"`
	RequiredLevel string `json:"required_level"`
}

type CareerSkillResponse struct {
	ID     string           `json:"career_id"`
	Name   string           `json:"name"`
	Desc   string           `json:"desc"`
	Skills []SkillsResponse `json:"skills"`
}

type CareerEditRequest struct {
	Name   string               `json:"name"`
	Desc   string               `json:"desc"`
	Skills []CareerSkillRequest `json:"skills"`
}
