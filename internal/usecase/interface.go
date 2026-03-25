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
}

type CareerRepository interface {
	FindAll(ctx context.Context) ([]entity.Career, error)
	FindById(ctx context.Context, id string) (*entity.Career, error)
	Update(ctx context.Context, career *entity.Career) error
	Delete(ctx context.Context, id string) error
	Save(ctx context.Context, career *entity.Career) error
	CreateCareerSkill(ctx context.Context, career *entity.Career, careerSkill []entity.CareerSkill) error
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

type SelfAssessmentRepository interface {
	CreateAssessmentSession(ctx context.Context, session *entity.UserCareerSession, skills []entity.SelfAssessmentSkill) error
}

type QuizRepository interface {
	GetSelfAssessmentSkillsBySession(ctx context.Context, sessionID string) ([]entity.SelfAssessmentSkill, error)
	GetRandomQuestionBySkillAndLevel(ctx context.Context, skillID string, level entity.LevelEnum) (*entity.Question, error)
	CreateQuizTransaction(ctx context.Context, quizSession *entity.QuizSession, quizAnswers []entity.QuizAnswer) error
}
