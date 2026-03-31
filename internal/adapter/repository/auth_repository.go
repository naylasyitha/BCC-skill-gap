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

// Save implements [usecase.AuthRepository].
func (a *authRepository) Save(ctx context.Context, user *entity.User) error {
	return a.db.WithContext(ctx).Create(user).Error
}

func (a *authRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := a.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update
func (a *authRepository) Update(ctx context.Context, user *entity.User) error {
	return a.db.WithContext(ctx).Save(user).Error
}

func (r *authRepository) UpdateRefreshToken(ctx context.Context, id string, token string) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		UpdateColumn("refresh_token", token).Error
}
