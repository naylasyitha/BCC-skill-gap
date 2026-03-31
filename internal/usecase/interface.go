package usecase

import (
	"context"
	"project-bcc/internal/entity"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	UpdateRefreshToken(ctx context.Context, id string, token string) error
}

type CareerRepository interface {
	FindAll(ctx context.Context) ([]entity.Career, error)
	FindById(ctx context.Context, id string) (*entity.Career, error)
	Update(ctx context.Context, career *entity.Career) error
	Delete(ctx context.Context, id string) error
	Save(ctx context.Context, career *entity.Career) error
	CreateCareerSkill(ctx context.Context, career *entity.Career, careerSkill []entity.CareerSkill) error
	UpdateCareerWithSkills(ctx context.Context, career *entity.Career, newSkills []entity.CareerSkill, updateSkills bool) error
}

type SkillRepository interface {
	FindAll(ctx context.Context) ([]entity.Skill, error)
	FindById(ctx context.Context, id string) (*entity.Skill, error)
	Update(ctx context.Context, skill *entity.Skill) error
	Delete(ctx context.Context, id string) error
	Save(ctx context.Context, skill *entity.Skill) error
}

type CareerSkillRepository interface {
	FindById(ctx context.Context, id string) (*entity.CareerSkill, error)
	Update(ctx context.Context, skill *entity.CareerSkill) error
	Delete(ctx context.Context, id string) error
	Save(ctx context.Context, skill *entity.CareerSkill) error
}

type CareerSessionRepository interface {
	Create(ctx context.Context, session *entity.UserCareerSession) error
	FindById(ctx context.Context, careerSessionId string) (*entity.UserCareerSession, error)
}

type SelfAssessmentRepository interface {
	CreateAssessmentSession(ctx context.Context, skills []entity.SelfAssessmentSkill) error
	UpdateStatus(ctx context.Context, careerSessionID string, status entity.StatusEnum) error
}

type QuizRepository interface {
	GetSelfAssessmentSkillsBySession(ctx context.Context, sessionID string) ([]entity.SelfAssessmentSkill, error)
	GetRandomQuestionBySkillAndLevel(ctx context.Context, skillID string, level entity.LevelEnum) (*entity.Question, error)
	CreateQuizTransaction(ctx context.Context, quizSession *entity.QuizSession, quizAnswers []entity.QuizAnswer) error
}

type QuestionRepository interface {
	Create(ctx context.Context, question *entity.Question) error
	FindAll(ctx context.Context) ([]entity.Question, error)
	FindById(ctx context.Context, id string) (*entity.Question, error)
	Update(ctx context.Context, question *entity.Question) error
	Delete(ctx context.Context, id string) error
}
