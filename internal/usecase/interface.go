package usecase

import (
	"context"
	"project-bcc/internal/entity"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
}

type CareerRepository interface {
	FindAll(ctx context.Context) ([]entity.Career, error)
	FindById(ctx context.Context, id string) (*entity.Career, error)
	Update(ctx context.Context, career *entity.Career) error
	Delete(ctx context.Context, id string) error
	Save(ctx context.Context, career *entity.Career) error
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
