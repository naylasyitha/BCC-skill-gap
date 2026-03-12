package repository

import (
	"context"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) usecase.AuthRepository {
	return &authRepository{db}
}

// FindByEmail implements [usecase.AuthRepository].
func (a *authRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := a.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername implements [usecase.AuthRepository].
func (a *authRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := a.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Save implements [usecase.AuthRepository].
func (a *authRepository) Save(ctx context.Context, user *entity.User) error {
	return a.db.WithContext(ctx).Create(user).Error
}
