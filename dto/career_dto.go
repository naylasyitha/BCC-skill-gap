package dto

type CareerCreateRequest struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc" binding:"required"`
}

type CareerResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type SkillsResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Priority      int    `json:"priority"`
	RequiredLevel string `json:"required_level"`
}

type CareerSkillResponse struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Desc   string           `json:"desc"`
	Skills []SkillsResponse `json:"skills"`
}

type CareerEditRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}
