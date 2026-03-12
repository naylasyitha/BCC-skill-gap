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
