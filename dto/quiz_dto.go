package dto

type QuizQuestionResponse struct {
	QuizAnswerID    string `json:"quiz_answer_id"`
	QuestionID      string `json:"question_id"`
	SkillID         string `json:"skill_id"`
	SkillName       string `json:"skill_name"`
	QuestionContent string `json:"question_content"`
	OptionA         string `json:"option_a"`
	OptionB         string `json:"option_b"`
	OptionC         string `json:"option_c"`
	OptionD         string `json:"option_d"`
}

type SkillFinalDetails struct {
	SkillID    string `json:"skill_id"`
	SkillName  string `json:"skill_name"`
	UserLevel  string `json:"user_level"`
	FinalLevel string `json:"final_user_level"`
	SkillScore int    `json:"skill_score"`
}

type StartQuizResponse struct {
	QuizSessionID   string                 `json:"quiz_session_id"`
	CareerSessionID string                 `json:"career_session_id"`
	Questions       []QuizQuestionResponse `json:"questions"`
}

type UpdateAnswerRequest struct {
	QuizAnswerID string `json:"quiz_answer_id" validate:"required"`
	UserAnswer   string `json:"user_answer" validate:"required,oneof=a b c d"`
}

type SubmitQuizResponse struct {
	QuizSessionID string              `json:"quiz_session_id"`
	TotalScore    int                 `json:"total_score"`
	SkillsResult  []SkillFinalDetails `json:"skills_result"`
}
