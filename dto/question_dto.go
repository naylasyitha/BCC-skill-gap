package dto

type QuestionCreateRequest struct {
	SkillID         string `json:"skill_id" binding:"required,uuid"`
	Level           string `json:"level" binding:"required"`
	QuestionContent string `json:"question_content" binding:"required"`
	OptionA         string `json:"option_a" binding:"required"`
	OptionB         string `json:"option_b" binding:"required"`
	OptionC         string `json:"option_c" binding:"required"`
	OptionD         string `json:"option_d" binding:"required"`
	Answer          string `json:"answer" binding:"required"`
	Explanation     string `json:"explanation"`
}

type QuestionResponse struct {
	ID              string `json:"question_id"`
	SkillID         string `json:"skill_id"`
	Level           string `json:"level"`
	QuestionContent string `json:"question_content"`
	OptionA         string `json:"option_a"`
	OptionB         string `json:"option_b"`
	OptionC         string `json:"option_c"`
	OptionD         string `json:"option_d"`
	Answer          string `json:"answer"`
	Explanation     string `json:"explanation"`
}

type QuestionUpdateRequest struct {
	SkillID         string `json:"skill_id"`
	Level           string `json:"level"`
	QuestionContent string `json:"question_content"`
	OptionA         string `json:"option_a"`
	OptionB         string `json:"option_b"`
	OptionC         string `json:"option_c"`
	OptionD         string `json:"option_d"`
	Answer          string `json:"answer"`
	Explanation     string `json:"explanation"`
}
