package repository

import (
	"context"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

	"gorm.io/gorm"
)

type careerSessionRepository struct {
	db *gorm.DB
}

// Create implements [usecase.CareerSessionRepository].

func NewCareerSessionRepository(db *gorm.DB) usecase.CareerSessionRepository {
	return &careerSessionRepository{db}
}

func (c *careerSessionRepository) Create(ctx context.Context, session *entity.UserCareerSession) error {
	return c.db.WithContext(ctx).Create(session).Error
}

func (c *careerSessionRepository) FindById(ctx context.Context, careerSessionId string) (*entity.UserCareerSession, error) {
	var careerSession entity.UserCareerSession
	err := c.db.WithContext(ctx).Preload("User").Preload("Career").Where("id = ?", careerSessionId).First(&careerSession).Error
	if err != nil {
		return nil, err
	}

	return &careerSession, nil
}
